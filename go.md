#  goLang's note
## 一:  goLang基础:
[goLangIDE服务器破解](http://intellij.mandroid.cn/)
[goLangIDE永久破解](http://www.jb51.net/softs/598316.html)
1. 数据类型: bool,int,string,指针,数组,struct,channel,函数,切片,接口,Map
	变量命名: 驼峰[区别局部变量和全局变量:首字母大小写]
	全局变量: 一个变量(常量,类型,函数) 在程序中都有一定的作用域.如果一个变量在函数外声明,则是全局变量,可以在整个包甚至外部包(被导入后)使用.不管你声明在哪个源文件里或在哪个源文件里调用该变量.
	局部变量: 它们的作用域只在函数体内,参数和返回值变量也是局部变量.函数体结束,该局部变量的生命周期也就结束了.一般情况下,局部变量的作用域可以通过代码块（用大括号括起来的部分）判断.
2. 可测试的Go代码(http://tabalt.net/blog/golang-testing/)
	Golang的测试代码位于某个包的源代码中名称以_test.go结尾的源文件里[sishen_test.go]
	测试代码包含[ 测试函数,测试辅助代码,示例函数 ];
	>测试函数: 有以Test开头的功能测试函数和以Benchmark开头的性能测试函数两种.

	>测试辅助: 代码是为测试函数服务的公共函数,初始化函数,测试数据等.

	>示例函数: 以Example开头的说明被测试函数用法的函数.

	大部分情况下,测试代码是作为某个包的一部分,意味着它可以访问包中不可导出的元素.
	但在有需要的时候（如避免循环依赖）也可以修改测试文件的包名,如package hello的测试文件,包名可以设为package hello_test.
		* Go标准库摘抄的 testing.T类型的常用方法的用法: 测试函数中的某条测试用例执行结果与预期不符时,
			调用t.Error()或t.Errorf()方法记录日志并标记测试失败
			使用t.Fatal()和t.Fatalf()方法,在某条测试用例失败后就跳出该测试函数
			使用t.Skip()和t.Skipf()方法,跳过某条测试用例的执行
			执行测试用例的过程中通过t.Log()和t.Logf()记录日志
			使用t.Parallel()标记需要并发执行的测试函数
		* 性能测试函数需要接收*testing.B类型的单一参数b,性能测试函数中需要循环b.N次调用被测函数.testing.B 类型用来管理测试时间和迭代运行次数,也支持和testing.T相同的方式管理测试状态和格式化的测试日志,不一样的是testing.B的日志总是会输出.
	    Go标准库摘抄的 testing.B类型的常用方法的用法：
			在函数中调用t.ReportAllocs(),启用内存使用分析
			通过 b.StopTimer(),b.ResetTimer(),b.StartTimer()来停止,重置,启动 时间经过和内存分配计数
			调用b.SetBytes()记录在一个操作中处理的字节数
			通过b.RunParallel()方法和 *testing.PB类型的Next()方法来并发执行被测对象
		* 测试辅助代码:(http://tabalt.net/blog/golang-testing/)
3. init函数初始化函数:
	init方法 是在任何package中都可以出现,但是建议 每个package中只包含一个init()函数比较好,容易理解.
	init在main函数执行之前执行.

4. 接口:
	Go语言数据结构的核心.是一些方法的集合.他指定了对象的行为: 如果它(任何数据类型)可以做这些事情,那么它就可以在这里使用.

5. := 结构不能使用在函数外.

6. 基本类型:
```
	Basic Types:
	bool
	string
	int  int8  int16  int32  int64
	uint uint8 uint16 uint32 uint64 uintptr
	byte // uint8 的别名
	rune // int32 的别名
	     // 代表一个Unicode码
	float32 float64
	complex64 complex128
```
7.  基本类型和引用类型:
	>基本类型(数组也是)因为是拷贝的值,并且在对他进行操作的时候,生成的也是新创建的值,所以这些类型在多线程里是安全的,我们不用担心一个线程的修改影响了另外一个线程的数据

	>引用类型和原始的基本类型恰恰相反,它的修改可以影响到任何引用到它的变量.在 Go 语言中,引用类型有 slice,map,chan,===接口,函数类型.

	>结构类型,要想修改结构实例的值,需要传递指针.[e. &p]
8.  数组(array)和切片(slice):
	>增加元素: append

	>遍历元素: for k,v := range mySlice{....}
9.  channel:
>管道是Go语言在语言级别上提供的goroutine间的**通讯方式**,我们可以使用channel在多个goroutine之间传递消息.channel是**进程内**的通讯方式,是不支持跨进程通信的,如果需要进程间通讯的话,可以使用Socket等网络方式
语法: var chanName chan ElemType  / ch := make(chan int, 100)<br/>
<br/>写入数据,在此需要注意：向管道中写入数据通常会导致程序阻塞,直到有其他goroutine从这个管道中读取数据ch<- value
<br/>读取数据,注意：如果管道中没有数据,那么从管道中读取数据会导致程序阻塞,直到有数据 value := <-ch
<br/>单向管道,var ch1 chan<- float64  / var ch2 <-chan int(只能读取int型数据)
<br/>关闭channel,直接调用close()即可 :  close(ch)
<br/>判断ch是否关闭,判断ok的值,如果是false,则说明已经关闭(关闭的话读取是不会阻塞的) x, ok := <-ch
<br/>Goroutine和channel是Go在“并发”方面两个核心feature.
<br/>channel 是 Goroutine 之间进行通信的一种方式(作用:同步 或传递信息).
<br/>channel 带缓冲的channel和不带缓冲的channel(区别:带缓冲的只有当写满了才会写阻塞,只有缓冲中一个也没有的时候才会读阻塞).
<br/>所谓的单向channel概念,其实只是对channel的一种使用限制.
<br/>单项channel初始化:ch3 := make(chan int) ; ch4 := <-chan int(ch3).
10. 协程(goroutine): 需探.
11. 类型方法:
12. new 和 make 区别:
>new 用在实例化结构体类型上比较能体现出价值.

	>make 的作用我深有体会,例如: 我只通过var m map[int]string;此时声明的是一个nil-map,故需要m = make(map[int]string)一下才能直接赋值.

	>一定要注意的是 嵌套类型时 ,必须对其make,否则会出错.例;

	```go
		var m map[int]map[int]string
		m = make(map[int]map[int]string)
		fmt.Println(m[2])
		v,ok := m[2]
		fmt.Println(v,ok)
		//此处的ok是false,说明m[2]没有被make定义,故需要检查ok
		if !ok{
			m[2] = make(map[int]string)
		}
```
13. goroutine 和 channel 是GO并发的两大基石,那么 接口 是Go语言编程中数据类型的关键.Go语言的实际编程中,几乎所有的数据结构都围绕接口展开,接口是Go语言中所有数据结构的核心.
14. 错误类型:Go内置了一个error类型,专门用来处理错误信息,
15. iota枚举:除非被显式设置为其它值或iota,每个const分组的第一个常量被默认设置为它的0值,第二及后续的常量被默认设置为它前面那个常量的值,如果前面那个常量的值是iota,则它也被设置为iota.
16. map读取的时候返回两个参数:
	```go
m := map[string]int{...}, m1,ok = m["..."]
```

17. 作用域
```go
	if x := 10;x > 1{
		fmt.Println("条件判断语句中允许声明一个变量,作用域为条件逻辑块",x)
	}else {
		fmt.Println("err")
	}
```

18. for配合range可以用于读取slice和map的数据.标签和break,continue,goto结合.
19. defer采用先进后出的方式.
20. 函数作为值,类型
	```go
	type testInt func(int) bool
	func isOdd(integer int) bool {
		if integer % 2 == 0{
			return false
		}else {
			return true
		}
	}
	func isEven(integer int) bool {
		if integer % 2 == 0{
			return true
		}else {
			return false
		}
	}
	func filter(slice []int,f testInt) []int {
		var result []int
		for _,value := range slice{
			if f(value){
				result = append(result,value)
			}
		}
		return result
	}
	func testFunc()  {
		slice := []int {1,2,3,4,5,7}
		fmt.Println("slice = " ,slice)
		odd := filter(slice,isOdd)
		fmt.Println("Odd elements of slice are:" , odd)
		even := filter(slice,isEven)
		fmt.Println("Even elements of slice are: ", even)
	}
```
21. panic:
	是一个内建函数,可以中断原有的控制流程,进入一个令人恐慌的流程中.当函数F调用panic,函数F的执行被中断,但是F中的延迟函数defer会正常执行,然后F返回到调用它的地方.在调用的地方,F的行为就像调用了panic.这一过程继续向上,直到发生panic的goroutine中所有调用的函数返回,此时程序退出.恐慌可以直接调用panic产生.也可以由运行时错误产生,例如访问越界的数组.
	```go
	func testPanic()  {
	var user = os.Getenv("USER")
		if user == "" {
			panic("no value for $USER")
		}
	}
```
22. recover:
	是一个内建的函数,可以让进入令人恐慌的流程中的goroutine恢复过来.recover仅在延迟函数defer中有效.在正常的执行过程中,调用recover会返回nil,并且没有其它任何效果.如果当前的goroutine陷入恐慌,调用recover可以捕获到panic的输入值,并且恢复正常的执行.
	```go
	func testRecover(f func()) (b bool) {
	defer func() {
		if x := recover();x != nil{
			b = true
			fmt.Println(x)
		}
	}()
	f() // 执行函数f,如果f中出现了panic,那么就可以恢复回来.
	return
	}
```
23. struct:
	>匿名字段(嵌入字段.或称为继承):

	>type定义的struct中的字段小写的作用域是当前包,大写的才可以在外部包中访问到.
24. method:
	>go中虽然没有class但是有method.

25. interface:(interface类型定义了一组方法,如果某个对象实现了某个接口的所有方法,则此对象就实现了此接口)
	>任意的类型都实现了空interface(我们这样定义：interface{}),也就是包含0个method的interface.
	<br/>如果我们定义了一个interface的变量,那么这个变量里面可以存实现这个interface的任意类型的对象
	<br/>验证某接口变量存储的是哪个类型,可以通过Comma-ok断言 或switch
	<br/>嵌入interface
26. 反射(reflect):
	<pre>
	使用reflect一般分为三步,要去反射一个类型的值(这些值都实现了空interface),
	* 首先需要把它转化成reflect对象(reflect.Type或者reflect.Value,根据不同情况调用不同的函数).这两种获取方式:
		①: t := reflect.TypeOf(i) // 得到类型的元数据,通过t我们能获取类型定义里面的所有元素
		②: v := reflect.ValueOf(i) // 得到实际的值,通过v我们获取存储在里面的值,还可以去改变值
	* 转化为reflect对象之后我们就可以进行一些操作了,也就是将reflect对象转化成相应的值
		tag := t.Elem().Field(0).Tag // 获取定义在struct里面的标签
		name:= v.Elem().Field(0).String() // 获取存储在第一个字段里面的值
	* 要想修改反射的字段的值,需要传址.
	* 反射的两个性质: https://studygolang.com/articles/2157
	* 反射定律: interface value to reflection object(接口类型值其实就是 谁实现了该接口)
				e: var i = interface{};var a int =10;i=a; 此时接口变量i中就包含属性(10,int)
				reflection object to interface value
	=====
	* 接口类型,一种特殊的数据类型,包括一组 方法集,同时,其内部维护两个属性：value 和 type.value指的是实现了接口类型的类型值
				;type则是对相应类型值的类型描述.
	* 反射定律: 1: 反射可以将“接口类型变量”转换为“反射类型对象(reflect.Type/reflect.Value:通过reflect.TypeOf()实现)”;
				2: 反射可以将“反射类型对象”转换为“接口类型变量”(根据一个 reflect.Value 类型的变量，我们可以使用 Interface 方法恢复其接口类型的值。事实上，这个方法会把 type 和 value 信息打包并填充到一个接口变量中)
				3: 如果要修改反射类型对象，其值必须是“addressable”
	</pre>
27. 测试Atom使用.
