#include "formatting/formatting.h"

#include <sstream>
#include <string>
#include <vector>

std::string JoinWithCommas(std::vector<std::string> parts) {
  std::stringstream out;
  if (parts.size() == 0) {
    return out.str();
  }
  out << parts[0];
  for (size_t i = 1; i < parts.size(); ++i) {
    out << ", " << parts[i];
  }
  return out.str();
}
