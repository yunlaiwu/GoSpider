# coding=utf-8


import sys

if __name__ == "__main__":
    try:
        action = sys.argv[1]
        id = sys.argv[2]
        title = sys.argv[3]
    except:
    	sys.stderr.write("\tpython " + sys.argv[0] + " action id title\n")
    	sys.exit(-1)

    mode = 'a' if action=='append' else 'w'
    fi = open('book_list.txt', mode)
    fi.write(id + '\t' + title + '\n')
    fi.close()