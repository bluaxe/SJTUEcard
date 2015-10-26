import socket
import urllib2

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

s.settimeout(5)

url = "http://ecard.sjtu.edu.cn/getCheckpic.action?rand=5475.172977894545"
resp = urllib2.urlopen(url)
data = resp.read()


s.connect(("127.0.0.1", 30196))
length = str(len(data))
print "length: "+length
s.send(length)
s.sendall(data)
res = s.recv(4)
s.close()

print res

