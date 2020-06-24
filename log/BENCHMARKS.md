Benchmarking this package, it should be clear that for simple logging, the standard library's log 
package performs better. Where this package starts to do better, however, is when there is a degree
of sophistication required in either the formatting of a log, or the information going into it (see 
for example the JSON formatted variants).

Another advantage this package has (if not in terms of pure performance) is the simplicity of use on
a call-by-call basis. To recreate similar functionality with the standard log package requires many
more lines of code. That's before getting into things like hooks.

The benchmarks show that if you're doing basic logging, the standard log is probably the way to go.
For any degree of sophistication and functionality, this package handles itself reasonably well, if 
not better than the standard log.
```
goos: linux
goarch: amd64
pkg: github.com/jackhascodes/kit/log
BenchmarkLog_Debug
BenchmarkLog_Debug-4                               	  376359	      2668 ns/op
BenchmarkLog_Log
BenchmarkLog_Log-4                                 	  477974	      2654 ns/op
BenchmarkStandardLog_Print
BenchmarkStandardLog_Print-4                       	 1007616	      1191 ns/op
BenchmarkLog_Debugf
BenchmarkLog_Debugf-4                              	  405159	      2797 ns/op
BenchmarkStandardLog_Printf
BenchmarkStandardLog_Printf-4                      	 1000000	      1260 ns/op
BenchmarkLog_Error
BenchmarkLog_Error-4                               	  240322	      5478 ns/op
BenchmarkStandardLog_Printf_withStackTrace
BenchmarkStandardLog_Printf_withStackTrace-4       	  267058	      4513 ns/op
BenchmarkLog_ErrorJson
BenchmarkLog_ErrorJson-4                           	  170494	      5907 ns/op
BenchmarkStandardLog_Printf_withStackTraceJson
BenchmarkStandardLog_Printf_withStackTraceJson-4   	  143268	      8280 ns/op
BenchmarkLog_DebugJson
BenchmarkLog_DebugJson-4                           	  365227	      3267 ns/op
BenchmarkStandardLog_PrintJson
BenchmarkStandardLog_PrintJson-4                   	  286036	      3921 ns/op
```