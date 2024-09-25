import tensorflow as tf
import tensorflow_hub as hub
import cv2
from . import video_util

# 加载预训练模型
detector = hub.load("https://www.kaggle.com/models/tensorflow/ssd-mobilenet-v2/TensorFlow2/ssd-mobilenet-v2/1")
# detector = hub.load("https://tfhub.dev/google/faster_rcnn/openimages_v4/inception_resnet_v2/1").signatures['default']


def detect_objects(frame):
    # 图像预处理
    img = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
    # img = tf.image.convert_image_dtype(img, tf.float32)
    img = tf.image.convert_image_dtype(img, tf.uint8)
    # img = tf.image.resize(img, [720, 640])
    img = tf.expand_dims(img, 0)
    # img = tf.cast(img, tf.uint8)

    # 进行物体检测
    results = detector(img)
    return results


async def run(video_path):
    frames, fps, width, height = video_util.load_video(video_path)
    key_frames = video_util.select_key_frames(frames, fps)
    response = []
    for frame in key_frames:
        results = detect_objects(video_util.resize_frame(frame))
        results = {key: results[key][0].numpy().tolist() for key in ["detection_boxes", "detection_classes", "detection_scores"]}
        response.append(results)
    return response
        # print(results)
        # detection_boxes = results["detection_boxes"][0]
        # detection_classes = results["detection_classes"][0]
        # detection_scores = results["detection_scores"][0]
        # for i in range(len(detection_boxes)):
        #     if detection_classes[i] > 0 and detection_scores[i] > 0.1:
        #         y_min, x_min, y_max, x_max = (detection_boxes[i][0] * height,
        #                                       detection_boxes[i][1] * width,
        #                                       detection_boxes[i][2] * height,
        #                                       detection_boxes[i][3] * width)
        #         cv2.rectangle(frame, (int(y_min), int(x_min)), (int(y_max), int(x_max)), (0, 255, 0), thickness=2)
        # cv2.imshow("Result", frame)
        # cv2.waitKey(0)