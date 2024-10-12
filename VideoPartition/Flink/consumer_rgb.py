# -*- coding: utf-8 -*-
from typing import Iterable
import time
from pyflink.datastream.connectors.file_system import FileSink, RollingPolicy, OutputFileConfig
from pyflink.datastream import StreamExecutionEnvironment, ProcessWindowFunction
from pyflink.datastream.connectors.kafka import FlinkKafkaConsumer
from pyflink.datastream.formats.json import JsonRowDeserializationSchema
from pyflink.datastream.functions import MapFunction, FilterFunction
from pyflink.common import Types, WatermarkStrategy, Time, Encoder
from pyflink.common.watermark_strategy import TimestampAssigner
from pyflink.datastream.window import EventTimeSessionWindows, TimeWindow

import cv2
import numpy as np


class OriginTimestampAssigner(TimestampAssigner):
    def extract_timestamp(self, value, record_timestamp) -> float:
        key_parts = value[0].split('_')
        return float(key_parts[3])


class MyTimestampAssigner(TimestampAssigner):
    def extract_timestamp(self, value, record_timestamp) -> float:
        return float(value[1])


class OrderProcessFunction(ProcessWindowFunction[tuple, tuple, str, TimeWindow]):
    def process(self,
                key: str,
                context: ProcessWindowFunction.Context[TimeWindow],
                elements: Iterable[tuple]) -> Iterable[tuple]:
        # 将元素按时间戳排序
        sorted_elements = sorted(elements, key=lambda x: x[1])

        for element in sorted_elements:
            yield element


class ToFrame(ProcessWindowFunction[tuple, tuple, str, TimeWindow]):
    def process(self,
                key: str,
                context: ProcessWindowFunction.Context[TimeWindow],
                elements: Iterable[tuple]) -> Iterable[tuple]:

        elements_iterator = iter(elements)
        first_element = next(elements_iterator, None)

        min_init_timestamp = float(first_element[0].split('_')[3])
        min_processing_timestamp = float(first_element[2])

        if first_element[0].split('_')[0] == 'trigger':
            return [(key, min_init_timestamp, None, min_processing_timestamp)]

        num_blocks = 10
        frame_width = first_element[1].shape[1]
        block_height = first_element[1].shape[0]
        frame_height = block_height * num_blocks

        restored_frame = np.zeros((frame_height, frame_width, 3), dtype=np.uint8)

        for element in elements:
            key_parts = element[0].split('_')

            if min_init_timestamp > float(key_parts[3]):
                min_init_timestamp = float(key_parts[3])
            if min_processing_timestamp > float(element[2]):
                min_processing_timestamp = float(element[2])

            index = int(key_parts[2])
            start_y = index * block_height
            end_y = (index + 1) * block_height
            start_x = 0
            end_x = frame_width

            # 将块复制到恢复的帧中
            restored_frame[start_y:end_y, start_x:end_x] = element[1]

        # frame_byte = cv2.imencode('.jpg', restored_frame)[1].tobytes()

        # return [(key, min_init_timestamp, frame_byte, min_processing_timestamp)]
        return [(key, min_init_timestamp, restored_frame, min_processing_timestamp)]


class GetMinProducingTimestamp(ProcessWindowFunction[tuple, tuple, str, TimeWindow]):
    def process(self,
                key: str,
                context: ProcessWindowFunction.Context[TimeWindow],
                elements: Iterable[tuple]) -> Iterable[tuple]:
        timestamp_list = []
        for element in elements:
            key_parts = element[0].split('_')
            if not key_parts[0] == "trigger":
                timestamp_list.append(float(key_parts[3]))
            else:
                if int(key_parts[1]) == 1:
                    return [(key, key_parts[2])]
                else:
                    return [('None', 'None')]
        return [(key, str(min(timestamp_list)))]


class AddTimeStamp(MapFunction):
    def map(self, x):
        return x[0], x[1], str(time.time())


class ToImage(MapFunction):
    def map(self, x):
        if x[1] is None:
            return x[0], None, x[2]

        img = cv2.imdecode(
            np.frombuffer(x[1], np.uint8), cv2.IMREAD_COLOR)
        return x[0], img, x[2]


class ToHist(MapFunction):
    def map(self, x):
        if x[2] is None:
            return x[0], x[1], -1, x[3]

        hist = cv2.calcHist([x[2]], [0, 1, 2], None, [8, 8, 8], [0, 256, 0, 256, 0, 256])
        cv2.normalize(hist, hist)

        return x[0], x[1], hist.flatten(), x[3]


class ToBytes(MapFunction):
    def map(self, x):
        if x[0].split('_')[0] == 'trigger':
            return x

        return x[0], x[1], x[2].tobytes(), x[3]


class D(MapFunction):
    def __init__(self):
        self.previous_frame = None
        self.previous_hist = None

    def map(self, x):
        current_frame = x[0].split('_')[1]

        if self.previous_hist is None or int(current_frame) < int(self.previous_frame):
            self.previous_frame = current_frame
            self.previous_hist = x[2]
            return x[1], x[3]  # 返回第一个元素的事件时间，用于计算延迟，第一个kafka端的事件时间，第二个为flink端的事件时间

        diff = cv2.compareHist(self.previous_hist, x[2], cv2.HISTCMP_BHATTACHARYYA)

        self.previous_hist = x[2]
        self.previous_frame = current_frame
        return x[0], x[1], 'pass', diff


class E(FilterFunction):
    def filter(self, x):
        return len(x) == 2 or float(x[3]) > 0.3


class ComputeLatency(MapFunction):
    def map(self, x):
        current_time = time.time()
        return f"{current_time - float(x[0])}, {x[1]}, {current_time}, {x[0]}"


def create_file_sink(base_path, prefix, suffix, part_size=1024 * 1024 * 128, rollover_interval=1 * 1000,
                     inactivity_interval=1 * 1000):
    return FileSink.for_row_format(
        base_path=base_path,
        encoder=Encoder.simple_string_encoder()
    ).with_output_file_config(
        OutputFileConfig.builder()
            .with_part_prefix(prefix)
            .with_part_suffix(suffix)
            .build()
    ).with_rolling_policy(
        RollingPolicy.default_rolling_policy(
            part_size=part_size,
            rollover_interval=rollover_interval,
            inactivity_interval=inactivity_interval
        )
    ).build()


if __name__ == '__main__':
    output_path = "/mnt/media/Ruiqi/pycharm_project/log/data"

    env = StreamExecutionEnvironment.get_execution_environment()
    env.set_parallelism(1)

    type_info = Types.ROW([Types.STRING(), Types.PRIMITIVE_ARRAY(Types.BYTE())])
    deserialization_schema = JsonRowDeserializationSchema.Builder().type_info(type_info).build()

    kafka_consumer = FlinkKafkaConsumer(
        topics='input-topic',
        deserialization_schema=deserialization_schema,
        properties={'bootstrap.servers': '128.110.96.30:9092', 'group.id': 'block-consumer'}
    )
    kafka_consumer.set_start_from_latest()

    data_stream = env.add_source(kafka_consumer)

    ds = data_stream.map(lambda x: (x[0], x[1])) \
        .map(AddTimeStamp())

    # ds.map(lambda x: (x[0], x[2])) \
    #   .print()

    origin_watermark_strategy = WatermarkStrategy.for_monotonous_timestamps() \
        .with_timestamp_assigner(OriginTimestampAssigner())

    frame_ds = ds.map(ToImage()) \
        .assign_timestamps_and_watermarks(origin_watermark_strategy) \
        .key_by(lambda x: "_".join(x[0].split('_')[:2]), key_type=Types.STRING()) \
        .window(EventTimeSessionWindows.with_gap(Time.milliseconds(5))) \
        .process(ToFrame())
    # .process(ToFrame(), Types.TUPLE([Types.STRING(), Types.DOUBLE(), Types.PRIMITIVE_ARRAY(Types.BYTE()), Types.DOUBLE()]))

    # frame_ds.print()

    # frame_ds.map(lambda x: (x[0], x[1], x[3])) \
    #   .print()

    watermark_strategy = WatermarkStrategy.for_monotonous_timestamps() \
        .with_timestamp_assigner(MyTimestampAssigner())

    ordered_ds = frame_ds.map(ToHist()) \
        .assign_timestamps_and_watermarks(watermark_strategy) \
        .key_by(lambda x: x[0].split('_')[0]) \
        .window(EventTimeSessionWindows.with_gap(Time.milliseconds(20))) \
        .process(OrderProcessFunction()) \
        .filter(lambda x: not (isinstance(x[2], int)))

    # .process(OrderProcessFunction(), Types.TUPLE([Types.STRING(), Types.DOUBLE(), Types.INT(), Types.DOUBLE()])) \

    # ordered_ds.map(lambda x: (x[0], x[1], x[3])) \
    #   .print()

    processed_ds = ordered_ds.map(D()) \
        .filter(E())

    # processed_ds.print()

    latency_ds = processed_ds.filter(lambda x: len(x) == 2) \
        .map(ComputeLatency(), output_type=Types.STRING())

    latency_ds.print()

    result_ds = processed_ds.filter(lambda x: not len(x) == 2) \
        .map(lambda x: f"{x[0]},{x[1]},{x[2]},{x[3]}", output_type=Types.STRING())

    result_ds.print()

    kafka_latency_ds = ds.assign_timestamps_and_watermarks(origin_watermark_strategy) \
        .key_by(lambda x: x[0].split('_')[0], key_type=Types.STRING()) \
        .window(EventTimeSessionWindows.with_gap(Time.milliseconds(5))) \
        .process(GetMinProducingTimestamp(), Types.TUPLE([Types.STRING(), Types.STRING()])) \
        .filter(lambda x: x[0] != 'None')

    kafka_latency_ds.print()

    # define the sink
    if output_path is not None:
        result_sink = create_file_sink(output_path, "result", ".txt")
        latency_sink = create_file_sink(output_path, "latency", ".txt")
        kafka_latency_sink = create_file_sink(output_path, "kafkaLatency", ".txt")

        result_ds.sink_to(result_sink)
        latency_ds.sink_to(latency_sink)
        kafka_latency_ds.sink_to(kafka_latency_sink)
    else:
        print("Use --output to specify output path.")

    env.execute("video2flink Consumer")
