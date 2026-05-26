
import os 
import sys
def Hello() :
    if len(sys.argv) != 2 :
        raise ValueError(
            f"Expected 2 args But Got {len(sys.argv)}\n"
            f"Usage python main.py priyankar"
        )
    print(f"{os.getenv("SOME_ENV_VAR")}")
    print(f"Hello This is From {sys.argv[1]}")

if __name__ == "__main__" : 
    Hello()