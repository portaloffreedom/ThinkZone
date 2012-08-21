#include <iostream>
#include "superstring.h"
#include "thinkzoneapp.h"
using namespace std;

int test_program()
{
    if (!SuperString::Test()) {
      cerr<<"#@# test sulle superstringe fallito! #@#\n";
      return 42;
    }
    
    return 0;
}

int main(int argc, char **argv) {
    cout << "Hello, world!" << endl;
//     return test_program();
    ThinkzoneApp a(argc,argv);
    return a.exec();
}
