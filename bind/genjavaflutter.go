package bind

func (g *JavaGen) flutterListFinder() {
	g.listFinder("call.args[0]", `result.fail("not found");`)
}

const (
	javaFlutterPreamble = `import io.flutter.plugin.common.MethodCall;
import io.flutter.plugin.common.MethodChannel;
import io.flutter.plugin.common.MethodChannel.MethodCallHandler;
import io.flutter.plugin.common.MethodChannel.Result;

`

)
