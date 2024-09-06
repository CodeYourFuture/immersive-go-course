+++
title="2. Common build actions"
+++

# 2. Common build actions

{{<note type="note" title="A note on languages">}}
Throughout this primer we will often use examples building and testing Java and C++ because they are languages with clear (and different) build steps, but the principles apply to most languages.


You do not need to know Java or C++ in order to follow this primer.


This is an interesting and common theme when working with build systems: we can often achieve a lot of things with a language without actually knowing many details of that language.


Tempting though it is to go down the rabbitholes of learning languages, we recommend you don't get too distracted with this unless you have a clear goal in doing so.
{{</note>}}

There are example projects for each of these languages in https://github.com/CodeYourFuture/immersive-go-course/tree/main/example-projects. Try to build them yourself on the command line to get a feel for each of these actions.

## Fetching external dependencies

When we depend on code that isn't part of our current project (perhaps because someone else wrote it, or because we re-use it across many projects), we need to fetch those dependencies onto our machine before we can build with it.

There are often two stages in this process: **resolving** and **fetching**. Most languages' tooling exposes this to users as one operation (e.g. `npm install`), rather than as two steps, but internally the tooling probably still does both.

### Resolving

The input to resolving is a list of dependencies. For instance, we may know we need "guava at exactly version `v33.2.1-jre`". Or we may know we need "guava at least version `v31.0-jre`". Or we may know we need "guava".

The output of resolving is a complete list of dependencies, including exact version numbers. For instance, a resolver may tell you: "You need guava at exactly version `v33.2.1-jre`, and guava needs junit at exactly version `5.10.2`, so you need that as well". Dependencies your project direectly needs are referred to as _direct dependencies_ and dependencies of your dependencies (or their dependencies) are referred to as _transitive dependencies_.

Different resolvers use different strategies for picking which versions you should use if your input doesn't have exact version requirements, but they always return exact versions. This means that if the inputs said "We need guava at least version `v31.0-jre`", the output of resolving will be "exactly `v31.0-jre`" or "exactly `v33.2.1-jre`" or some other exact version.

### Fetching

Fetching typically just involves downloading a file and putting it on disk in the right place for your project to use it (possibly unzipping it on the way). Where "on disk in the right place" is depends on the build system being used. Using `npm` in JavaScript, for instance, it's a directory named after the dependency, in a directory called `node_modules` in the same directory as your project's `package.json`.

## Compiling

Compiling is the process of taking source code (intended for humans to read/write) and converting it into a format better suited for computers to consume and run.

### C++ (a "native" language)

In C++, compiling takes a `.cc` file and generates a `.o` file. This `.o` file contains actual machine code that an operating system knows how to run directly. There are different compilers for C++: common ones are called `clang++`, `g++`, or `c++`.

A typical command line to compile C++ with some dependency is:

```bash
clang++ formatting/formatting.cpp -I "." -c -o formatting/formatting.o
```

The `-o` flag specifies where the generated object file should be placed.

The `-I` flag points at a directory containing `.h` files, telling the compiler "if you need to include a `.h` file, here's a place to look for them". It may be specified multiple times. `.h` files are header files, which give the compiler enough information about code you depend on for your code to be able to use it (generally type definitions and function signatures).

### Java (a "JIT" language)

In Java, compiling takes a `.java` file and generates a `.class` file. This `.class` file contains JVM bytecode which a `java` process (called a JVM) knows how to run. When running, the JVM converts the bytecode to machine code that an operating system knows how to run directly.

A typical command line to compile Java with some dependency is:

```bash
javac -cp third_party/guava-33.2.1-jre.jar com/example/fmt/Formatting.java
```

The `-cp` flag specifies directories or jar files (a jar file is a zip file with an expected structure inside) where `.class` files you depend on can be found.

### JavaScript (an "interpreted" language)

JavaScript does not have a standard compile step - it is an interpreted language, which means that when you run `node` (or load a file in a browser), the JavaScript engine reads that file as text, line-by-line, and executes it. But JavaScript development often involves steps like bundling (e.g. with `webpack`, which often creates a new `.js` file from an existing one) which in some ways resemble compilation.

### Comparisons

These three compilation modes offer different trade-offs:
* Native machine code is very fast to run (because it's directly run by the operating system), and during the compile step optimisations may be applied to speed up the code further (at the cost of taking even longer to compile).
* Interpreted languages are much slower to run (because they can't be optimised in advance), but don't require a compile step in order to run (so are often faster to iterate on when writing code).
* Java compilation does some compilation ahead of time (which takes some time, but typically less time than compiling native machine code), but by deferring some decisions until runtime, the JVM can observe how code is actually used with real-world data and perform optimisations in response to that (this is called Just-In-Time compilation or JIT).

## Linking

Linking is the process of taking several pieces of compiled code (e.g. object files of machine code) and combining them into one executable.

### C++

Before linking, an object file may reference some symbols (e.g. function names) that it just knows the name of. We can run `nm` on an object file and see what symbols it contains and which ones it knows it needs but doesn't know about. In the output of `nm`, the first column is the address of the symbol in the object file (if it's present), the second describes its type (where `U` means "undefined" - not in this object file, but something in this file needs it), and the third column is the symbol's name.

<details>
<summary>

Click to expand an `nm` invocation from the example project's `main.o`.
</summary>

```console
% nm main.o
0000000000002a3c s GCC_except_table0
0000000000002af0 s GCC_except_table103
0000000000002b00 s GCC_except_table110
0000000000002b10 s GCC_except_table111
0000000000002b20 s GCC_except_table115
0000000000002b30 s GCC_except_table116
0000000000002b64 s GCC_except_table118
0000000000002a68 s GCC_except_table12
0000000000002b7c s GCC_except_table138
0000000000002b8c s GCC_except_table143
0000000000002b9c s GCC_except_table147
0000000000002bb0 s GCC_except_table150
0000000000002bc0 s GCC_except_table153
0000000000002a78 s GCC_except_table21
0000000000002a8c s GCC_except_table28
0000000000002a50 s GCC_except_table4
0000000000002a9c s GCC_except_table41
0000000000002aac s GCC_except_table57
0000000000002abc s GCC_except_table76
0000000000002ad0 s GCC_except_table88
0000000000002ae0 s GCC_except_table98
0000000000002a28 s __GLOBAL__sub_I_main.cc
                 U __Unwind_Resume
                 U __Z14JoinWithCommasNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEEE
0000000000003fe0 b __ZN9constantsL5namesE
00000000000007f0 T __ZNKSt16initializer_listINSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEEE3endB7v160006Ev
0000000000000680 T __ZNKSt16initializer_listINSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEEE4sizeB7v160006Ev
00000000000007d8 T __ZNKSt16initializer_listINSt3__112basic_stringIcNS0_11char_traitsIcEENS0_9allocatorIcEEEEE5beginB7v160006Ev
0000000000002124 T __ZNKSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE13__get_pointerB7v160006Ev
0000000000002818 T __ZNKSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE15__get_long_sizeB7v160006Ev
0000000000002840 T __ZNKSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE16__get_short_sizeB7v160006Ev
00000000000021b0 T __ZNKSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE18__get_long_pointerB7v160006Ev
00000000000021d8 T __ZNKSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE19__get_short_pointerB7v160006Ev
0000000000002068 T __ZNKSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE4dataB7v160006Ev
00000000000027c4 T __ZNKSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE4sizeB7v160006Ev
0000000000002178 T __ZNKSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE9__is_longB7v160006Ev
0000000000001c7c T __ZNKSt3__113basic_ostreamIcNS_11char_traitsIcEEE6sentrycvbB7v160006Ev
00000000000015a0 T __ZNKSt3__116reverse_iteratorIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEEE4baseB7v160006Ev
000000000000163c T __ZNKSt3__116reverse_iteratorIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEEEdeB7v160006Ev
0000000000001618 T __ZNKSt3__116reverse_iteratorIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEEEptB7v160006Ev
0000000000002200 T __ZNKSt3__117__compressed_pairINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE5__repES5_E5firstB7v160006Ev
0000000000000fd4 T __ZNKSt3__117__compressed_pairIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE5firstB7v160006Ev
0000000000000ca0 T __ZNKSt3__117__compressed_pairIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE6secondB7v160006Ev
0000000000001f88 T __ZNKSt3__119ostreambuf_iteratorIcNS_11char_traitsIcEEE6failedB7v160006Ev
0000000000002224 T __ZNKSt3__122__compressed_pair_elemINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE5__repELi0ELb0EE5__getB7v160006Ev
0000000000000cc4 T __ZNKSt3__122__compressed_pair_elemINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEELi1ELb1EE5__getB7v160006Ev
0000000000000ff8 T __ZNKSt3__122__compressed_pair_elemIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEELi0ELb0EE5__getB7v160006Ev
00000000000013dc T __ZNKSt3__129_AllocatorDestroyRangeReverseINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEPS6_EclB7v160006Ev
00000000000023bc T __ZNKSt3__15ctypeIcE5widenB7v160006Ec
0000000000000c54 T __ZNKSt3__16__lessImmEclB7v160006ERKmS3_
                 U __ZNKSt3__16locale9use_facetERNS0_2idE
0000000000000aac T __ZNKSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE14__annotate_newB7v160006Em
0000000000001798 T __ZNKSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE17__annotate_deleteB7v160006Ev
0000000000000a00 T __ZNKSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE20__throw_length_errorB7v160006Ev
0000000000000f10 T __ZNKSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE31__annotate_contiguous_containerB7v160006EPKvSA_SA_SA_
0000000000000f30 T __ZNKSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE4dataB7v160006Ev
00000000000018d8 T __ZNKSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE4sizeB7v160006Ev
0000000000000bc4 T __ZNKSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE7__allocB7v160006Ev
0000000000000f58 T __ZNKSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE8capacityB7v160006Ev
00000000000009a0 T __ZNKSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE8max_sizeEv
0000000000000fac T __ZNKSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE9__end_capB7v160006Ev
0000000000001f0c T __ZNKSt3__18ios_base5flagsB7v160006Ev
00000000000022c4 T __ZNKSt3__18ios_base5rdbufB7v160006Ev
0000000000001fd8 T __ZNKSt3__18ios_base5widthB7v160006Ev
                 U __ZNKSt3__18ios_base6getlocEv
0000000000000c88 T __ZNKSt3__19allocatorINS_12basic_stringIcNS_11char_traitsIcEENS0_IcEEEEE8max_sizeB7v160006Ev
0000000000001f24 T __ZNKSt3__19basic_iosIcNS_11char_traitsIcEEE4fillB7v160006Ev
00000000000022a0 T __ZNKSt3__19basic_iosIcNS_11char_traitsIcEEE5rdbufB7v160006Ev
000000000000230c T __ZNKSt3__19basic_iosIcNS_11char_traitsIcEEE5widenB7v160006Ec
                 U __ZNSt11logic_errorC2EPKc
0000000000000d44 T __ZNSt12length_errorC1B7v160006EPKc
0000000000000d78 T __ZNSt12length_errorC2B7v160006EPKc
                 U __ZNSt12length_errorD1Ev
                 U __ZNSt20bad_array_new_lengthC1Ev
                 U __ZNSt20bad_array_new_lengthD1Ev
00000000000022dc T __ZNSt3__111char_traitsIcE11eq_int_typeEii
0000000000002304 T __ZNSt3__111char_traitsIcE3eofEv
000000000000036c T __ZNSt3__111char_traitsIcE6lengthEPKc
0000000000002110 T __ZNSt3__112__to_addressB7v160006IKcEEPT_S3_
0000000000000f98 T __ZNSt3__112__to_addressB7v160006INS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEEEEPT_S8_
000000000000155c T __ZNSt3__112__to_addressB7v160006INS_16reverse_iteratorIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEEEEvEENS_5decayIDTclsr19__to_address_helperIT_EE6__callclsr3stdE7declvalIRKSB_EEEEE4typeESD_
                 U __ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE6__initEPKcm
                 U __ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE6__initEmc
000000000000202c T __ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEC1B7v160006Emc
0000000000000000 T __ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEC1B7v160006IDnEEPKc
                 U __ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEC1ERKS5_
00000000000020bc T __ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEC2B7v160006Emc
00000000000002d0 T __ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEC2B7v160006IDnEEPKc
                 U __ZNSt3__112basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEED1Ev
                 U __ZNSt3__113basic_ostreamIcNS_11char_traitsIcEEE3putEc
                 U __ZNSt3__113basic_ostreamIcNS_11char_traitsIcEEE5flushEv
                 U __ZNSt3__113basic_ostreamIcNS_11char_traitsIcEEE6sentryC1ERS3_
                 U __ZNSt3__113basic_ostreamIcNS_11char_traitsIcEEE6sentryD1Ev
000000000000024c T __ZNSt3__113basic_ostreamIcNS_11char_traitsIcEEElsB7v160006EPFRS3_S4_E
0000000000000bec T __ZNSt3__114numeric_limitsIlE3maxB7v160006Ev
0000000000002238 T __ZNSt3__114pointer_traitsIPKcE10pointer_toB7v160006ERS1_
0000000000001ff0 T __ZNSt3__115basic_streambufIcNS_11char_traitsIcEEE5sputnB7v160006EPKcl
0000000000000914 T __ZNSt3__116__non_trivial_ifILb1ENS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEEC2B7v160006Ev
0000000000000458 T __ZNSt3__116__non_trivial_ifILb1ENS_9allocatorIcEEEC2B7v160006Ev
0000000000001c98 T __ZNSt3__116__pad_and_outputB7v160006IcNS_11char_traitsIcEEEENS_19ostreambuf_iteratorIT_T0_EES6_PKS4_S8_S8_RNS_8ios_baseES4_
00000000000018a4 T __ZNSt3__116allocator_traitsINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEE10deallocateB7v160006ERS7_PS6_m
0000000000002538 T __ZNSt3__116allocator_traitsINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEE37select_on_container_copy_constructionB7v160006IS7_vvEES7_RKS7_
0000000000001530 T __ZNSt3__116allocator_traitsINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEE7destroyB7v160006IS6_vEEvRS7_PT_
0000000000000ba0 T __ZNSt3__116allocator_traitsINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEE8max_sizeB7v160006IS7_vEEmRKS7_
0000000000001244 T __ZNSt3__116allocator_traitsINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEE9constructB7v160006IS6_JRKS6_EvEEvRS7_PT_DpOT0_
0000000000002760 T __ZNSt3__116allocator_traitsINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEE9constructB7v160006IS6_JRS6_EvEEvRS7_PT_DpOT0_
00000000000014b4 T __ZNSt3__116reverse_iteratorIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEEEC1B7v160006ES7_
0000000000001664 T __ZNSt3__116reverse_iteratorIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEEEC2B7v160006ES7_
0000000000001580 T __ZNSt3__116reverse_iteratorIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEEEppB7v160006Ev
0000000000000330 T __ZNSt3__117__compressed_pairINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE5__repES5_EC1B7v160006INS_18__default_init_tagESA_EEOT_OT0_
00000000000003b0 T __ZNSt3__117__compressed_pairINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE5__repES5_EC2B7v160006INS_18__default_init_tagESA_EEOT_OT0_
0000000000000ed8 T __ZNSt3__117__compressed_pairIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE5firstB7v160006Ev
0000000000000ea0 T __ZNSt3__117__compressed_pairIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE6secondB7v160006Ev
00000000000005c0 T __ZNSt3__117__compressed_pairIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEEC1B7v160006IDnNS_18__default_init_tagEEEOT_OT0_
0000000000002548 T __ZNSt3__117__compressed_pairIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEEC1B7v160006IDnS8_EEOT_OT0_
0000000000000860 T __ZNSt3__117__compressed_pairIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEEC2B7v160006IDnNS_18__default_init_tagEEEOT_OT0_
0000000000002610 T __ZNSt3__117__compressed_pairIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEEC2B7v160006IDnS8_EEOT_OT0_
0000000000000e54 T __ZNSt3__117__libcpp_allocateB7v160006Emm
000000000000046c T __ZNSt3__118__constexpr_strlenB7v160006EPKc
000000000000186c T __ZNSt3__118__debug_db_erase_cB7v160006INS_6vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS5_IS7_EEEEEEvPT_
0000000000000a1c T __ZNSt3__119__allocate_at_leastB7v160006INS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEEENS_19__allocation_resultINS_16allocator_traitsIT_E7pointerEEERSA_m
0000000000001444 T __ZNSt3__119__allocator_destroyB7v160006INS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEENS_16reverse_iteratorIPS6_EESA_EEvRT_T0_T1_
00000000000003a0 T __ZNSt3__119__debug_db_insert_cB7v160006INS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEEEEvPT_
0000000000000670 T __ZNSt3__119__debug_db_insert_cB7v160006INS_6vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS5_IS7_EEEEEEvPT_
00000000000019d0 T __ZNSt3__119__libcpp_deallocateB7v160006EPvmm
00000000000015e0 T __ZNSt3__119__to_address_helperINS_16reverse_iteratorIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEEEEvE6__callB7v160006ERKS9_
0000000000001ed8 T __ZNSt3__119ostreambuf_iteratorIcNS_11char_traitsIcEEEC1B7v160006ERNS_13basic_ostreamIcS2_EE
000000000000224c T __ZNSt3__119ostreambuf_iteratorIcNS_11char_traitsIcEEEC2B7v160006ERNS_13basic_ostreamIcS2_EE
0000000000000ce0 T __ZNSt3__120__throw_length_errorB7v160006EPKc
0000000000000e7c T __ZNSt3__121__libcpp_operator_newB7v160006IJmEEEPvDpT_
00000000000003ec T __ZNSt3__122__compressed_pair_elemINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEE5__repELi0ELb0EEC2B7v160006ENS_18__default_init_tagE
0000000000000ec4 T __ZNSt3__122__compressed_pair_elemINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEELi1ELb1EE5__getB7v160006Ev
00000000000008bc T __ZNSt3__122__compressed_pair_elemINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEELi1ELb1EEC2B7v160006ENS_18__default_init_tagE
0000000000002654 T __ZNSt3__122__compressed_pair_elemINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEELi1ELb1EEC2B7v160006IS7_vEEOT_
0000000000000400 T __ZNSt3__122__compressed_pair_elemINS_9allocatorIcEELi1ELb1EEC2B7v160006ENS_18__default_init_tagE
0000000000000efc T __ZNSt3__122__compressed_pair_elemIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEELi0ELb0EE5__getB7v160006Ev
00000000000008a0 T __ZNSt3__122__compressed_pair_elemIPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEELi0ELb0EEC2B7v160006IDnvEEOT_
00000000000011b4 T __ZNSt3__122__make_exception_guardB7v160006INS_29_AllocatorDestroyRangeReverseINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS2_IcEEEEEEPS7_EEEENS_28__exception_guard_exceptionsIT_EESC_
00000000000005fc T __ZNSt3__122__make_exception_guardB7v160006INS_6vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS5_IS7_EEE16__destroy_vectorEEENS_28__exception_guard_exceptionsIT_EESC_
0000000000000cd8 T __ZNSt3__123__libcpp_numeric_limitsIlLb1EE3maxB7v160006Ev
0000000000001a28 T __ZNSt3__124__libcpp_operator_deleteB7v160006IJPvEEEvDpT_
0000000000001a90 T __ZNSt3__124__put_character_sequenceB7v160006IcNS_11char_traitsIcEEEERNS_13basic_ostreamIT_T0_EES7_PKS4_m
0000000000001a00 T __ZNSt3__127__do_deallocate_handle_sizeB7v160006IJEEEvPvmDpT_
0000000000001278 T __ZNSt3__128__exception_guard_exceptionsINS_29_AllocatorDestroyRangeReverseINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS2_IcEEEEEEPS7_EEE10__completeB7v160006Ev
00000000000012c0 T __ZNSt3__128__exception_guard_exceptionsINS_29_AllocatorDestroyRangeReverseINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS2_IcEEEEEEPS7_EEEC1B7v160006ESA_
00000000000012f4 T __ZNSt3__128__exception_guard_exceptionsINS_29_AllocatorDestroyRangeReverseINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS2_IcEEEEEEPS7_EEEC2B7v160006ESA_
0000000000001294 T __ZNSt3__128__exception_guard_exceptionsINS_29_AllocatorDestroyRangeReverseINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS2_IcEEEEEEPS7_EEED1B7v160006Ev
000000000000138c T __ZNSt3__128__exception_guard_exceptionsINS_29_AllocatorDestroyRangeReverseINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS2_IcEEEEEEPS7_EEED2B7v160006Ev
0000000000000818 T __ZNSt3__128__exception_guard_exceptionsINS_6vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS5_IS7_EEE16__destroy_vectorEE10__completeB7v160006Ev
0000000000000928 T __ZNSt3__128__exception_guard_exceptionsINS_6vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS5_IS7_EEE16__destroy_vectorEEC1B7v160006ESA_
000000000000095c T __ZNSt3__128__exception_guard_exceptionsINS_6vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS5_IS7_EEE16__destroy_vectorEEC2B7v160006ESA_
0000000000000834 T __ZNSt3__128__exception_guard_exceptionsINS_6vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS5_IS7_EEE16__destroy_vectorEED1B7v160006Ev
00000000000016ac T __ZNSt3__128__exception_guard_exceptionsINS_6vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS5_IS7_EEE16__destroy_vectorEED2B7v160006Ev
0000000000001200 T __ZNSt3__129_AllocatorDestroyRangeReverseINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEPS6_EC1B7v160006ERS7_RS8_SB_
0000000000001324 T __ZNSt3__129_AllocatorDestroyRangeReverseINS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEPS6_EC2B7v160006ERS7_RS8_SB_
0000000000001048 T __ZNSt3__130__uninitialized_allocator_copyB7v160006INS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEPKS6_S9_PS6_EET2_RT_T0_T1_SB_
000000000000266c T __ZNSt3__130__uninitialized_allocator_copyB7v160006INS_9allocatorINS_12basic_stringIcNS_11char_traitsIcEENS1_IcEEEEEEPS6_S8_S8_EET2_RT_T0_T1_S9_
0000000000000b74 T __ZNSt3__13minB7v160006ImEERKT_S3_S3_
0000000000000c00 T __ZNSt3__13minB7v160006ImNS_6__lessImmEEEERKT_S5_S5_T0_
                 U __ZNSt3__14coutE
0000000000000278 T __ZNSt3__14endlB7v160006IcNS_11char_traitsIcEEEERNS_13basic_ostreamIT_T0_EES7_
                 U __ZNSt3__15ctypeIcE2idE
                 U __ZNSt3__16localeD1Ev
0000000000000698 T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE11__vallocateB7v160006Em
000000000000063c T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE16__destroy_vectorC1ERS8_
0000000000000980 T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE16__destroy_vectorC2ERS8_
00000000000016fc T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE16__destroy_vectorclB7v160006Ev
000000000000074c T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE18__construct_at_endIPKS6_Li0EEEvT_SC_m
0000000000002584 T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE18__construct_at_endIPS6_Li0EEEvT_SB_m
000000000000100c T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE21_ConstructTransactionC1B7v160006ERS8_m
0000000000001168 T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE21_ConstructTransactionC2B7v160006ERS8_m
000000000000113c T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE21_ConstructTransactionD1B7v160006Ev
000000000000168c T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE21_ConstructTransactionD2B7v160006Ev
0000000000001900 T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE22__base_destruct_at_endB7v160006EPS6_
0000000000000a5c T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE7__allocB7v160006Ev
000000000000187c T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE7__clearB7v160006Ev
0000000000000a84 T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEE9__end_capB7v160006Ev
0000000000000034 T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEEC1B7v160006ESt16initializer_listIS6_E
0000000000000218 T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEEC1ERKS8_
000000000000049c T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEEC2B7v160006ESt16initializer_listIS6_E
0000000000002428 T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEEC2ERKS8_
0000000000000070 T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEED1B7v160006Ev
0000000000001a4c T __ZNSt3__16vectorINS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEENS4_IS6_EEED2B7v160006Ev
                 U __ZNSt3__18ios_base33__set_badbit_and_consider_rethrowEv
                 U __ZNSt3__18ios_base5clearEj
0000000000002090 T __ZNSt3__18ios_base5widthB7v160006El
00000000000023f4 T __ZNSt3__18ios_base8setstateB7v160006Ej
000000000000198c T __ZNSt3__19allocatorINS_12basic_stringIcNS_11char_traitsIcEENS0_IcEEEEE10deallocateB7v160006EPS5_m
00000000000015b8 T __ZNSt3__19allocatorINS_12basic_stringIcNS_11char_traitsIcEENS0_IcEEEEE7destroyB7v160006EPS5_
0000000000000dc4 T __ZNSt3__19allocatorINS_12basic_stringIcNS_11char_traitsIcEENS0_IcEEEEE8allocateB7v160006Em
000000000000135c T __ZNSt3__19allocatorINS_12basic_stringIcNS_11char_traitsIcEENS0_IcEEEEE9constructB7v160006IS5_JRKS5_EEEvPT_DpOT0_
0000000000002794 T __ZNSt3__19allocatorINS_12basic_stringIcNS_11char_traitsIcEENS0_IcEEEEE9constructB7v160006IS5_JRS5_EEEvPT_DpOT0_
00000000000008e8 T __ZNSt3__19allocatorINS_12basic_stringIcNS_11char_traitsIcEENS0_IcEEEEEC2B7v160006Ev
000000000000042c T __ZNSt3__19allocatorIcEC2B7v160006Ev
0000000000001fac T __ZNSt3__19basic_iosIcNS_11char_traitsIcEEE8setstateB7v160006Ej
0000000000002390 T __ZNSt3__19use_facetB7v160006INS_5ctypeIcEEEERKT_RKNS_6localeE
0000000000000184 T __ZNSt3__1lsB7v160006INS_11char_traitsIcEEEERNS_13basic_ostreamIcT_EES6_PKc
00000000000001cc T __ZNSt3__1lsB7v160006IcNS_11char_traitsIcEENS_9allocatorIcEEEERNS_13basic_ostreamIT_T0_EES9_RKNS_12basic_stringIS6_S7_T1_EE
00000000000014e8 T __ZNSt3__1neB7v160006IPNS_12basic_stringIcNS_11char_traitsIcEENS_9allocatorIcEEEES7_EEbRKNS_16reverse_iteratorIT_EERKNS8_IT0_EE
0000000000000e20 T __ZSt28__throw_bad_array_new_lengthB7v160006v
                 U __ZSt9terminatev
                 U __ZTISt12length_error
                 U __ZTISt20bad_array_new_length
                 U __ZTVSt12length_error
                 U __ZdlPv
                 U __Znwm
0000000000000490 T ___clang_call_terminate
                 U ___cxa_allocate_exception
                 U ___cxa_atexit
                 U ___cxa_begin_catch
                 U ___cxa_end_catch
                 U ___cxa_free_exception
                 U ___cxa_throw
0000000000002874 s ___cxx_global_var_init
                 U ___dso_handle
                 U ___gxx_personality_v0
                 U ___stack_chk_fail
                 U ___stack_chk_guard
000000000000009c T _main
                 U _strlen
0000000000002bd4 s l_.str
0000000000002bdb s l_.str.1
0000000000002be1 s l_.str.2
0000000000002be8 s l_.str.3
0000000000002bea s l_.str.4
0000000000000000 t ltmp0
0000000000002874 s ltmp1
0000000000002a3c s ltmp2
0000000000003fe0 b ltmp3
0000000000002bd4 s ltmp4
0000000000002bf8 s ltmp5
0000000000002c00 s ltmp6
```
</details>

The linker's job is to make sure it can find a definition for each symbol any object file references, and to stitch them all together so that when one function tries to call another one, the other function can be found in the executable and called.

## Packaging

When preparing software to be deployed somewhere, it often needs to be packaged in some format.

That may be a directory with a certain structure, a zip file with a certain structure, or something more complicated. This may be tied specifically to the language being built, or may be dictated by the environment where it's being deployed.

## Pre-processing / generating code

Sometimes code may be very repetitive, or may be derived from some data.

For instance, if we are using protobuf files to define types and functions that may be used, we need to run the protocol buffer compiler to generate the per-language bindings from the .proto files.

A different example: if we have a CSV of categories, and need to populate some constants or a hashmap with them, we can generate code to do this at compile-time rather than needing to parse the CSV at run-time. This may be beneficial for a number of reasons: we don't need to distribute the CSV with our code, we don't need to worry about the speed of access of a file from disk, we don't need to handle possible errors from reading, etc.

These are examples which may take some input (e.g. a .proto or .csv file) and generate some code, which we then need to compile. This generation needs to be done before we can compile the code, but once we have generated the code it can be treated like any other code that's in our project.

## Post-processing

Sometimes after we have produced some output (e.g. by compiling some code), we have extra steps we want to perform (perhaps before packaging it, or perhaps to analyse the output in some way).

Some examples:
* There may be some optimisations (e.g. minifying code to reduce binary size, or doing [profile-guided optimisations](https://en.wikipedia.org/wiki/Profile-guided_optimization)) we want to perform which edit a compiled output to improve it in some way.
* We may want to scan the compiled output for secrets which shouldn't be compiled in the code and which should instead be read from a file/environment - if the secrets are detected, we want to give an error indicating someone should fix this.

## Testing

### Output

In many ways, testing is one of the most unusual kind of build actions.

People may have many different signals they're looking for from tests, depending on their context.

A lot of the time, people just want to know the answer to a boolean question: Did all the tests pass?

If some tests failed, they probably want to know which ones, and what error messages were reported.

A lot of testing frameworks will also produce structured output (often an XML file) which other tools can consume. For example, an IDE may display the results of the tests. Or if every time tests are run the XML is uploaded to some service, it can track and present data like what are your slowest tests, or what tests fail the most often.

### Reliability

Tests are also unusual because they are often less consistent than other kinds of build actions. If you compile some code, you should always end up with the same output (or the same error, if it couldn't compile!). If you package some files in a zip file, you should end up with the same zip file. If you run tests, sometimes you don't get the same output.

Sometimes this is intentional. When tests run, they may log information, and that may include things like timestamps, to help you understand what happened when the test was running.

Sometimes this is accidental. Some tests may depend on how fast something happens, and if your computer was too busy (so the test ran slowly) it may fail. Some tests only work if the computer they're running on is set to a particular time zone, so you may get different results than if one of your colleagues runs the test.

Sometimes we want to run tests lots of times, e.g. to find out whether they sometimes randomly fail, so that we can detect problems and fix them.

This is not usually the case with other kinds of build actions. Compiling is consistent. Packaging is consistent. Good tests are consistent, but not all tests are good tests. Often, this means build systems need to treat test actions differently from other actions.

### How are tests actually run?

Tests are run via a program called a test harness.

In some languages, the test harness is a standalone program which is given the tests to run. This is quite common in interpreted languages, e.g. in JavaScript, `jest` tests are run by running a program called `jest`, and that program will look for test files it should run.

In other languages, when you want to run tests, a program gets compiled (where perhaps a `main` function is generated based on what tests exist), and then to run the tests you run that program. This means every time a test is changed or added, some code may need to be re-generated, and a whole program needs to be compiled.

The important thing to take away from this is that running tests is really just running a program. The tests are considered to have passed if that program exits 0. The tests are considered to have failed if the program exits non-zero. The test program will write output to stdout and stderr, and maybe produce structured output (e.g. an XML file) somewhere too. But running tests is generally just: running a program, and interpreting its output.

## Observations

What do all of these kinds of build action have in common?

* They typically have some input files (maybe source code, maybe a CSV file, maybe a compiled test program).
* They typically have some output files (maybe object files, maybe source code, maybe test results).
* They are generally performed by running a process.
* They can succeed or fail, and this is generally signalled by the exit code of a process.

Build actions may also have ordering constraints. If we're generating code, we need to do that before we can compile it. If we're running tests, we need to compile the tests before we can run them. If we're running tests which include generated code, we need to generate the code before we can compile the code before we can run the tests.
