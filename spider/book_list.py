# coding=utf-8


import sys

if __name__ == "__main__":
    try:
        res = sys.argv[1]
        action = sys.argv[2]
        id = sys.argv[3]
        title = sys.argv[4]
    except:
    	sys.stderr.write("\tpython " + sys.argv[0] + " action id title\n")
    	sys.exit(-1)

    filename = 'booklist.txt' if res == 'book' else 'movielist.txt'
    mode = 'a' if action=='append' else 'w'
    fi = open(filename, mode)
    fi.write(id + '\t' + title + '\n')
    fi.close()