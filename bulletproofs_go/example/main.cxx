#include "libbulletproofs.h"
#include <iostream>

using namespace std;

int main() {
  cout << "C++ says: calling Go dynamic lib.." << endl;
  GoSetuoGeneric(30, 12);
  cout << "C++ says: got total as "  << endl;
}
