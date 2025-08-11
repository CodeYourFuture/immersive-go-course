package com.example.fmt;

import com.google.common.base.Joiner;

public class Formatting {
  public static String joinWithCommas(Iterable<String> names) {
    return Joiner.on(", ").join(names);
  }
}
