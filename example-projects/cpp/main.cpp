#include <iostream>

#include "formatting/formatting.h"
#include "gen/constants.h"

int main(void) {
  std::cout << "Hello " << JoinWithCommas(constants::names) << "!" << std::endl;
  return 0;
}
