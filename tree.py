import os
from directory_tree import DisplayTree 

# pip install directory_tree
if __name__ == '__main__':
    
    curr_dir = os.path.dirname(os.path.abspath(__file__))
    print("Current Directory:", curr_dir)
    
    DisplayTree(curr_dir)
