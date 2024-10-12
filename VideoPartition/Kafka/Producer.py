import argparse
from kafka import KafkaProducer
import sys
import json
import cv2
import base64
import time
import threading


class MyEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, bytes):
            return base64.b64encode(obj).decode('utf-8')
        return json.JSONEncoder.default(self, obj)


def open_video(path):
    cap = cv2.VideoCapture(path)
    if (cap.isOpened()):
        print("open video file: " + path + " success!")
        return cap
    else:
        print("open video file: " + path + " failed!")
        sys.exit()


def write_to_kafka_thread(path, producer, timestamps):
    write_to_kafka(path, producer, timestamps)


def write_to_kafka(path, producer, timestamps, frame_interval=90, parallelism=1, num_blocks=10):
    filename = path.split("/")[-1]
    cap = open_video(path)
    fps = cap.get(cv2.CAP_PROP_FPS)
    total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
    print("Total frames:", total_frames)

    current_second_timestamp = None
    for frame_idx in range(0, total_frames, frame_interval):
        print("Extract frame: " + path + "_" + str(frame_idx))
        cap.set(cv2.CAP_PROP_POS_FRAMES, frame_idx)
        ret, frame = cap.read()
        if not ret:
            continue
        current_frame_second = int(cap.get(cv2.CAP_PROP_POS_MSEC))

        block_width = frame.shape[1]
        block_height = frame.shape[0] // num_blocks

        for i in range(num_blocks):
            current_second_timestamp = time.time()

            start_y = i * block_height
            end_y = (i + 1) * block_height
            start_x = 0
            end_x = block_width
            # 从帧中提取块
            block = frame[start_y:end_y, start_x:end_x]

            # 转换为字节数组
            block_byte = cv2.imencode('.jpg', block)[1].tobytes()
            key = filename + "_" + str(current_frame_second) + "_" + str(i) + "_" + str(current_second_timestamp)
            print(key)

            record = {"f0": key, "f1": block_byte}
            future = producer.send(topic='input-topic', value=record).get(timeout=10)

    timestamps.append(current_second_timestamp)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Process some files.')
    parser.add_argument('-p', '--path', nargs='+', help='The path to the file.')
    args = parser.parse_args()

    producer = KafkaProducer(
        bootstrap_servers=['128.110.96.30:9092'],
        value_serializer=lambda m: json.dumps(m, cls=MyEncoder).encode(),
        api_version=(0, 10)
    )

    if args.path:
        threads = []
        timestamps = []
        for path in args.path:
            # 创建并启动线程
            thread = threading.Thread(target=write_to_kafka_thread, args=(path, producer, timestamps))
            thread.start()
            threads.append(thread)

        # 等待所有线程完成
        for thread in threads:
            thread.join()

        # 获取最后一个时间戳
        if timestamps:
            final_timestamp = max(timestamps)
            key = "trigger_1_" + str(time.time()) + "_" + str(final_timestamp + 21)
            record = {"f0": key, "f1": None}
            future = producer.send(topic='input-topic', value=record).get(timeout=10)

            key = "trigger_2_2_" + str(final_timestamp + 31)
            record = {"f0": key, "f1": None}
            future = producer.send(topic='input-topic', value=record).get(timeout=10)
    else:
        print("No paths provided.")
