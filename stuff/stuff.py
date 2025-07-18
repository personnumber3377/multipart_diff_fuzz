
import os

DIR = "corpus/"
OUT = "output/"
def readfile(f):
	fh = open(DIR+f, "r")
	data = fh.read()
	fh.close()
	return data

def savefile(data, i):
	fh = open(OUT+str(i), "w")
	fh.write(data)
	fh.close()

def main():
	files = os.listdir(DIR)
	i = 0
	for f in files:
		data = readfile(f)
		# boundary=
		if "boundary=" in data:
			boundary = data[data.index("boundary=")+len("boundary="):]
			boundary = boundary[:boundary.index("\n")]
			print("Boundary: "+str(boundary))
			data = data.replace(boundary, "--RubyBoundary") # Just do something like this maybe???
			data = data.replace("boundary=--RubyBoundary", "boundary=RubyBoundary")
			savefile(data, i)
			i += 1
	return
if __name__=="__main__":
	main()
	exit()
