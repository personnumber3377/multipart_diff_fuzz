
import sys

def usage():
        print("Usage: python "+str(sys.argv[0])+" GOINPUTFILE OUTPUTBINARYFILE")
        return

def save_binary_to_file(filename, binary):
        fh = open(filename, "wb")
        fh.write(binary)
        fh.close()

def file_readlines(filename): # Read all lines including line separator and return them as a list.
        fh = open(filename, "r")
        lines = fh.readlines()
        fh.close()
        return lines # Return

def main():
        if len(sys.argv) != 3:
                usage()
                exit(1)
        infile = sys.argv[1]
        outfile = sys.argv[2]
        bytes_line = file_readlines(infile)[1] # Get the second file from the line.
        sep = "[]byte(\""
        assert bytes_line.startswith(sep) # Should be the thing...
        bytes_line = bytes_line[len(sep):-3] # Cutoff the thing...
        # print(bytes_line)
        bytes_line = bytes(bytes_line, "ascii") # Inefficient, but I don't really care...
        bytes_line_literal = bytes_line.decode('unicode_escape') # "unescape" the string.
        bytes_line_literal = bytes(bytes_line_literal, "ascii")
        # print(bytes_line_literal)
        save_binary_to_file(outfile, bytes_line_literal)
        return



if __name__=="__main__":
        main()
        exit()
