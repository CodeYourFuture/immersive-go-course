package com.example;

import com.example.fmt.Formatting;
import com.example.gen.Constants;

class Main {
  public static void main(String[] args) {
    System.out.printf("Hello %s!%n", Formatting.joinWithCommas(Constants.names()));
  }
}
