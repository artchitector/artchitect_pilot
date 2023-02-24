import io
import os
import re

from PIL import Image
from flask import Flask, Response, request

app = Flask(__name__)


@app.route('/painting', methods=['POST'])
def painting():
    print('tags: ' + request.form['tags'])
    print('seed: ' + request.form['seed'])
    print('width: ' + request.form['width'])
    print('height: ' + request.form['height'])
    print('steps: ' + request.form['steps'])
    print('version: ' + request.form['version'])

    filename = getPaintingFromInvokeAIFilename(request.form['version'])

    im = Image.open(filename)
    print(im.format, im.size, im.mode)

    img_byte_arr = io.BytesIO()
    im.save(img_byte_arr, format="PNG")

    return Response(img_byte_arr.getvalue(), content_type="image/png")


def getPaintingFromInvokeAIFilename(version):
    prepareFileForInvokeAI(version)


    filename = None

    if version.find("v1") == 0:
        cmd = '/home/artchitector/invoke-ai/invokeai_v2.2.5/invokeai/.venv/bin/python ' \
              + '/home/artchitector/invoke-ai/invokeai_v2.2.5/invokeai/.venv/bin/invoke.py ' \
            + '--from_file "/home/artchitector/invoke-ai/invokeai_v2.2.5/invokeai/list.txt"'
        pattern = re.compile(".*(\/home\/artchitector\/invoke-ai\/invokeai_v2.2.5\/invokeai\/outputs\/[0-9\.]+png).*")
    elif version.find("v2") == 0:
        cmd = 'INVOKEAI_ROOT=/home/artchitector/invoke-ai/invokeai_v2.3.0/ ' \
              + '/home/artchitector/invoke-ai/invokeai_v2.3.0/.venv/bin/python /home/artchitector/invoke-ai/invokeai_v2.3.0/.venv/bin/invoke.py ' \
              + '--from_file "/home/artchitector/invoke-ai/invokeai_v2.3.0/list.txt"'
        pattern = re.compile(".*(\/home\/artchitector\/invoke-ai\/invokeai_v2.3.0\/outputs\/[0-9\.]+png).*")
    else:
        print("unknown version")
        exit(1)

    ret = os.popen(cmd)
    lines = ret.readlines()
    for line in lines:
        match = pattern.match(line)
        if match is not None:
            filename = match.groups()[0]
            print(f"Found filename: {filename}")

    if filename is not None:
        return filename
    else:
        raise Exception("filename not found")


def prepareFileForInvokeAI(version):
    tags = request.form['tags']
    seed = request.form['seed']
    width = request.form['width']
    height = request.form['height']
    steps = request.form['steps']
    upscale = request.form['upscale']

    lines = []
    lines.append(f'{tags} -S{seed} -W{width} -H{height} -s{steps} -U{upscale}')
    if version.find("v1") == 0:
        filename = "/home/artchitector/invoke-ai/invokeai_v2.2.5/invokeai/list.txt"
    elif version.find("v2") == 0:
        filename = "/home/artchitector/invoke-ai/invokeai_v2.3.0/list.txt"
    else:
        print("unknown version")
        exit(1)
    with open(filename, "w") as text_file:
        for line in lines:
            text_file.write(line)
    text_file.close()


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8083, debug=True)
