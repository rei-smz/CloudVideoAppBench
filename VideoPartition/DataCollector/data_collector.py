import subprocess
import re
import time
import os
import signal
import pandas as pd
import numpy as np
import yaml


def start_kafka(host, files):
    command = '"/users/jliu1721/miniconda3/envs/pyflink/bin/python3.8 /mnt/media/Ruiqi/pycharm_project/src/producer_optimize.py -p '
    file_paths = ''

    for file in files:
        file_paths = file_paths + '/mnt/media/Ruiqi/pycharm_project/sample/' + file + '.mp4 '
    command = command + file_paths + '"'
    remote_command = f"ssh -t {host} {command}"

    print('kafka running...')

    process = subprocess.Popen(remote_command, shell=True, preexec_fn=os.setsid,
                               stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

    process.wait()
    print("kafka producer completed!")


def start_collectl(host, prefix):
    file_path = None
    if host is not None:
        command = f'"/mnt/media/Ruiqi/pycharm_project/log/collectl_log.sh scCM /mnt/media/Ruiqi/pycharm_project/log/data/{prefix}_{host}.txt"'
        remote_command = f"ssh -t {host} {command}"
        process = subprocess.Popen(remote_command, shell=True, preexec_fn=os.setsid, stdout=subprocess.PIPE,
                                   stderr=subprocess.PIPE)
        file_path = f'/mnt/media/Ruiqi/pycharm_project/log/data/{prefix}_{host}.txt'
    else:
        command = f'/mnt/media/Ruiqi/pycharm_project/log/collectl_log.sh scCM /mnt/media/Ruiqi/pycharm_project/log/data/{prefix}.txt'
        process = subprocess.Popen(command, shell=True, preexec_fn=os.setsid, stdout=subprocess.PIPE,
                                   stderr=subprocess.PIPE)
        file_path = f'/mnt/media/Ruiqi/pycharm_project/log/data/{prefix}.txt'
    if host is not None:
        print(f'hosts:{host} collectl running...')
    else:
        print('local collectl running...')
    return process, file_path


def kill_process(process):
    print(f'kill process: {process.pid}')
    os.killpg(os.getpgid(process.pid), signal.SIGINT)
    process.wait()


def check_directory(path):
    try:
        # 获取指定目录下的所有文件和子目录
        items = os.listdir(path)
        for item in items:
            full_path = os.path.join(path, item)
            # 检查是否是目录且以"2024"开头
            if os.path.isdir(full_path) and item.startswith("2024"):
                return full_path
        return None
    except FileNotFoundError:
        return None
    except Exception as e:
        print(f"An error occurred: {e}")
        return None


def check_files_nums(path):
    try:
        # 获取指定目录下的所有文件和子目录
        items = os.listdir(path)
        files = []

        # 遍历所有项目，找到文件
        for item in items:
            full_path = os.path.join(path, item)
            if os.path.isfile(full_path):
                files.append(full_path)

        # 过滤文件名开头分别为 kafkaLatency, latency, result 的文件
        filtered_files = []
        for file in files:
            file_name = os.path.basename(file)
            if file_name.startswith("kafkaLatency"):
                filtered_files.append(file)
            elif file_name.startswith("latency"):
                filtered_files.append(file)
            elif file_name.startswith("result"):
                filtered_files.append(file)

        # 如果找到的文件数量等于3，返回这些文件的绝对路径
        if len(filtered_files) == 3:
            return filtered_files[:3]
        else:
            return None
    except FileNotFoundError:
        print(f"The directory {path} does not exist.")
        return None
    except Exception as e:
        print(f"An error occurred: {e}")
        return None


def process_result(path):
    data = pd.read_csv(path, header=None)
    first_column = data[0]

    result_dic = {}
    for item in first_column:
        parts = item.split('_')
        filename = parts[0].split('-')[0]
        frame_second = int(parts[1])

        if filename not in result_dic:
            result_dic[filename] = set()

        result_dic[filename].add(frame_second)

    file_path = '/mnt/media/Ruiqi/pycharm_project/log/data/result.txt'
    with open(file_path, 'a', encoding='utf-8') as file:
        for key in result_dic:
            file.write(f'*****{key}*****\n')
            frames = sorted(result_dic[key])
            for second in frames:
                file.write(f'{second}\n')


def process_kafka_time(path):
    with open(path, 'r', encoding='utf-8') as file:
        end = None
        start_time_list = []
        for init_line in file:
            line = init_line.strip().strip('()')
            parts = line.split(',')
            filename = parts[0]
            timestamp = parts[1]

            if filename.startswith('trigger'):
                end = float(timestamp)
            elif filename.startswith('sample'):
                start_time_list.append(float(timestamp))
        start = max(start_time_list)
    return start, end


def process_latency(path):
    with open(path, 'r', encoding='utf-8') as file:
        latency_list = []
        for init_line in file:
            line = init_line.strip()
            latency_list.append(float(line.split(', ')[0]))

    return max(latency_list)


def process_flink_time(path):
    with open(path, 'r', encoding='utf-8') as file:
        start_list = []
        end_list = []
        for init_line in file:
            line = init_line.strip()
            parts = line.split(', ')

            start_time = parts[1]
            end_time = parts[2]

            start_list.append(float(start_time))
            end_list.append(float(end_time))

    return max(start_list), min(end_list)


def save_cpu_usage(filename, start, end, data):
    output = ','.join(map(str, data))
    with open(f'/mnt/media/Ruiqi/pycharm_project/log/data/{filename}', 'a', encoding='utf-8') as file:
        file.write(f'start time:{start}\n')
        file.write(f'end time:{end}\n')
        file.write(f'data:{output}\n')


def clean_up(timestamp):
    print('cleaning up...')
    original_path = '/mnt/media/Ruiqi/pycharm_project/log/data'
    new_path = f'/mnt/media/Ruiqi/pycharm_project/log/data_{timestamp}'

    command = f'mv {original_path} {new_path} && mkdir -p /mnt/media/Ruiqi/pycharm_project/log/data'
    subprocess.run(command, shell=True, check=True)
    print('cleaning up completed!')


def process_data(files, sample_files):
    print('processing data...')

    flink_data = '/mnt/media/Ruiqi/pycharm_project/log/data/flink.txt'
    kafka_data = '/mnt/media/Ruiqi/pycharm_project/log/data/kafka_jliu1721@node-3.txt'
    output_file = '/mnt/media/Ruiqi/pycharm_project/log/data/output.yaml'
    config_file = '/mnt/media/Ruiqi/pycharm_project/src/config.yaml'

    flink_cpu_usage = None
    kafka_cpu_usage = None
    latency = None
    for file_path in files:
        filename = file_path.split('/')[-1]
        if filename.startswith('latency'):
            latency = process_latency(file_path)
            flink_start, flink_end = process_flink_time(file_path)
            flink_cpu_usage = compute_avg_cpu_usage(flink_data, flink_start, flink_end)
            save_cpu_usage('flink_cpu_usage.txt', flink_start, flink_end, flink_cpu_usage)
        elif filename.startswith('result'):
            process_result(file_path)
        elif filename.startswith('kafkaLatency'):
            kafka_start, kafka_end = process_kafka_time(file_path)
            kafka_cpu_usage = compute_avg_cpu_usage(kafka_data, kafka_start, kafka_end)
            save_cpu_usage('kafka_cpu_usage.txt', kafka_start, kafka_end, kafka_cpu_usage)

    with open(config_file, 'r') as file:
        config = yaml.safe_load(file)
    frame_interval = config['producer']['frame.interval']
    timestamp = time.time()
    with open(output_file, 'a', encoding='utf-8') as file:
        file.write('config:\n')
        file.write(f'frame.interval:{frame_interval}\n')
        file.write('flink.config:standalone\n')
        file.write('kafka.config:standalone\n')
        input_file = ','.join(map(str, sample_files))
        file.write(f'input.file:{input_file}\n')
        file.write('label:\n')
        file.write(f'timestamp:{timestamp}\n')
        file.write('result:\n')
        file.write(f'latency:{latency}\n')
        file.write(f'flink.avg.cpu:{str(np.mean(flink_cpu_usage))}\n')
        file.write(f'kafka.avg.cpu:{str(np.mean(kafka_cpu_usage))}\n')

    print('data process completed!')
    print(f'see: /mnt/media/Ruiqi/pycharm_project/log/data_{timestamp}/output.yaml')

    clean_up(timestamp)


def scp_file2local(host, remote_filepath):
    print('moving remote file to local: ...')
    command = f'scp {host}:{remote_filepath} /mnt/media/Ruiqi/pycharm_project/log/data'
    subprocess.run(command, shell=True, check=True)

    # 删除远程服务器上的文件
    rm_command = f'rm {remote_filepath}'
    remote_command = f"ssh {host} {rm_command}"
    subprocess.run(remote_command, shell=True, check=True)


def extract_timestamp(line):
    match = re.search(r'\((\d+\.\d+)\)', line)
    if match:
        return match.group(1)
    return None


def compute_avg_cpu_usage(filepath, start, end):
    usage_list = []
    with open(filepath, 'r') as f:
        lines = iter(f)
        for line in lines:
            if '### RECORD' in line:
                timestamp = float(extract_timestamp(line))
                if start < timestamp < end:
                    for _ in range(3):
                        next(lines)
                    fourth_line = next(lines)
                    words = fourth_line.split()
                    value = int(words[9])
                    usage_list.append(100 - value)

    return usage_list


if __name__ == '__main__':
    host = "jliu1721@node-3"
    sample_files = input('input video filename, different video using space to distinguish\n').split(' ')

    kafka_collectl, kafka_data_file = start_collectl(host, 'kafka')
    local_collectl, local_data_file = start_collectl(None, 'flink')
    time.sleep(3)

    start_kafka(host, sample_files)

    print('waiting generate data...')

    directory_to_check = "/mnt/media/Ruiqi/pycharm_project/log/data"
    full_path = None
    while True:
        full_path = check_directory(directory_to_check)
        if full_path:
            break
        time.sleep(1)

    files = None
    while True:
        files = check_files_nums(full_path)
        if files:
            print('data generated success!')
            break
        else:
            time.sleep(1)

    kill_process(kafka_collectl)
    kill_process(local_collectl)

    scp_file2local(host, kafka_data_file)

    process_data(files, sample_files)
