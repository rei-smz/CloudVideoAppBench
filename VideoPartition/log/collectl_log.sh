

# example: scCM file.txt

if [ $# -ne 2 ]; then
    echo "Usage: $0 <option> <output_file>"
    exit 1
fi

# 定义一个函数来处理 kill 信号
cleanup() {
    # 在脚本被 kill 前执行一些命令
    echo "Received kill signal. Cleaning up..."
    #gzip -d $2-*.raw.gz > "unzip_$2"
    #rm $2-*.raw.gz
    exit 1
}

trap 'cleanup "$1" "$2"' SIGINT SIGTERM

#collectl -$1 -i 0.2 -f $2
collectl -$1 -i 0.2 >> $2
