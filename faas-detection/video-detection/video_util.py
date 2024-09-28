import cv2


def load_video(video_path):
    # 使用 OpenCV 加载视频
    cap = cv2.VideoCapture(video_path)
    frames = []
    fps = int(cap.get(cv2.CAP_PROP_FPS))  # 获取视频的帧率
    width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))  # 获取视频的宽度
    height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))  # 获取视频的高度
    try:
        while True:
            ret, frame = cap.read()
            if not ret:
                break
            frames.append(frame)
    finally:
        cap.release()
    return frames, fps, width, height


def resize_frame(frame, scale=0.5):
    return cv2.resize(frame, (int(frame.shape[1] * scale), int(frame.shape[0] * scale)), interpolation=cv2.INTER_AREA)


def select_key_frames(frames, fps):
    # 简单选择每秒的第一帧作为关键帧
    key_frames = [frames[i] for i in range(0, len(frames), fps)]
    return key_frames

def preprocess(video_path):
    frames, fps, width, height = load_video(video_path)
    key_frames = select_key_frames(frames, fps)
    return width, height, key_frames
