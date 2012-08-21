#include <iostream>
#include "superstring.h"
using namespace std;

int test_program()
{
    if (!SuperString::Test()) {
      cerr<<"#@# test sulle superstringe fallito! #@#\n";
      return 42;
    }
}

int main(int argc, char **argv) {
    cout << "Hello, world!" << endl;
    return test_program();
}
