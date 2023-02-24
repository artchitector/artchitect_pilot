import io
import os
import re
import subprocess

from flask import Flask, render_template, Response, abort, request
from io import BytesIO
from PIL import Image

app = Flask(__name__)


@app.route('/painting', methods=['POST'])
def painting():
    print(request.form['tags'])
    print(request.form['seed'])
    print(request.form['width'])
    print(request.form['height'])
    print(request.form['steps'])

    filename = getPaintingFromInvokeAIFilename()

    im = Image.open(filename)
    print(im.format, im.size, im.mode)

    img_byte_arr = io.BytesIO()
    im.save(img_byte_arr, format="PNG")

    return Response(img_byte_arr.getvalue(), content_type="image/png")


def getPaintingFromInvokeAIFilename():
    prepareFileForInvokeAI()

    pattern = re.compile(".*(\/home\/artchitector\/invoke-ai\/invokeai_v2.2.5\/invokeai\/outputs\/[0-9\.]+png).*")
    filename = None

    cmd = '/home/artchitector/invoke-ai/invokeai_v2.2.5/invokeai/.venv/bin/python /home/artchitector/invoke-ai/invokeai_v2.2.5/invokeai/.venv/bin/invoke.py --from_file "/home/artchitector/invoke-ai/invokeai_v2.2.5/invokeai/list.txt"'
    print(cmd)
    ret = os.popen(
        cmd
        )
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


def prepareFileForInvokeAI():
    tags = request.form['tags']
    seed = request.form['seed']
    width = request.form['width']
    height = request.form['height']
    steps = request.form['steps']
    upscale = request.form['upscale']

    lines = []
    lines.append(f'{tags} -S{seed} -W{width} -H{height} -s{steps} -U{upscale}')
    with open("/home/artchitector/invoke-ai/invokeai_v2.2.5/invokeai/list.txt", "w") as text_file:
        for line in lines:
            text_file.write(line)
    text_file.close()


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8083, debug=True)
