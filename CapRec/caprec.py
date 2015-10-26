from PIL import Image
from thread import start_new_thread
import urllib2
import StringIO
import socket
import sys

box = (9, 12)
offset = [2, 3]
th = 100

def neis(pos):
	x = pos[0]
	y = pos[1]
	nei = []
	for i in range(x-1, x+2):
		for j in range (y-1, y+2):
			nei.append((i,j))
	return nei

def similarity(a, b):
	total = box[0] * box[1]
	cnt = 0
	ap = a.load()
	bp = b.load()
	for x in range(box[0]):
		for y in range(box[1]):
			if ap[x,y] == bp[x,y]:
				cnt +=1
	return float(cnt)/ total

def rec(roi):
	num = 0
	maximum = 0
	for i in range(10):
		mod = Image.open("model/"+str(i)+".jpg")		
		sim = similarity(roi, mod)
		#print i,sim
		if sim > maximum:
			maximum = sim
			num = i
	return num


#print img.size
def thresh(img):
	pixels = img.load()
	for x in range(img.size[0]):
		for y in range(img.size[1]):
			if pixels[x,y] > th :
				pixels[(x,y)] = 255
			else:
				pixels[x,y] = 0
	return img
	
	#img.save("2.jpg")


def dec_noise(img):
	nonoise = Image.new("L", img.size)

	pix_no = nonoise.load()
	pixels = img.load()
	
	for x in range(1, img.size[0]-1):
		for y in range(1, img.size[1]-1):
			pix_no[x,y] = pixels[x,y]
			if pixels[x,y] < th:
				cnt = 0
				for pos in neis((x,y)):
					if pixels[pos[0], pos[1]] > th:
						cnt +=1
				#print (x,y), cnt
				if cnt >= 7:
					pix_no[x,y] = 255

	#nonoise.save("3.jpg")
	return nonoise

def recognize(data):
	img = Image.open(StringIO.StringIO(data))
	img = img.convert("L")
	img.save("client.jpg")
	
	img = thresh(img)
	img = dec_noise(img)

	boxs = []
	off = [offset[0], offset[1]]
	for i in range(4):
		boxs.append((off[0], off[1], off[0]+box[0], off[1]+box[1]))
		off[0] += box[0]
	
	res = ""
	i=1
	for b in boxs:
		roi = img.crop(b)	
		#roi.save("c"+str(i)+".jpg")
		i+=1
		num = rec(roi)
		res +=str(num)
	return res

def handler(conn):
	length = conn.recv(4)
	print "len:"+length
	data = conn.recv(int(length))
	
	res = recognize(data)
	print "code is :"+res
	conn.send(res)
	conn.close()

def server_start(port):
	try:
		s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
		s.bind(("0.0.0.0", port))

		s.listen(30)
		print "Start server at port:"+str(port)
		while True:
			conn, addr = s.accept()
			print addr[0] + ":"+str(addr[1])+" connected."
			start_new_thread(handler, (conn,))
	except KeyboardInterrupt :
		pass
	except Exception as e:
		print e
	finally:
		s.close()
	
	
'''
url = "http://ecard.sjtu.edu.cn/getCheckpic.action?rand=5475.172977894545"
resp = urllib2.urlopen(url)
data = resp.read()
res = recognize(data)
print res
'''

server_start(30196)

