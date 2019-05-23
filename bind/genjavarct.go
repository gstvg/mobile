package bind

import (
	"go/types"
	"golang.org/x/mobile/internal/importers/java"
)

func (g *JavaGen) GenRCTJavaClass(idx int) error {


	return nil
}

func (g *JavaGen) GenRCT() error {

	if g.isRn {
		g.Printf("/*")

		if len(g.constants) > 0 {
			g.Printf("@Override")
			g.Printf("public HashMap<String, Object> getConstants() {")
			g.Indent()
			g.Printf("final Map<String, Object> constants = new HashMap<>();")
			for _, constant := range g.constants {
				g.Printf(`constants.put("%s", this.%s);`, constant.Name(), constant.Name())
			}
			g.Outdent()
			g.Printf("}\n")
		}
		g.Printf("@Override\n")
		g.Printf("public String getName() {\n")
		g.Indent()
		g.Printf(`return "%s";`, g.className())
		g.Outdent()
		g.Printf("}\n")
		g.Printf("*/")
	}


	return nil
}

func (g *JavaGen) genRCTStruct(s structInfo) {

	if g.isRn {
		if jinf != nil {
			if jinf.extends == nil {
				jinf.extends = &java.Class{Name: "ReactContextBaseJavaModule"}
			} else {
			}
		} else {
		}
	}

	if g.isFlutter {
		impls = append(impls, "MethodCallHandler")
	}


	if g.isFlutter || g.isRn {

		g.Printf("/*")
		if g.isFlutter {
			g.Printf("public final class %ss implements MethodCallHandler {\n", n)
		} else {
			g.Printf("public final class %ss extends ReactContextBaseJavaModule {\n", n)
		}
		g.Indent()
		g.Printf("\nprivate static HashMap<String, %s> list = new HashMap<>();", n)

		if g.isFlutter {

			g.Printf("public static void registerWith(Registrar registrar) {\n")
			g.Indent()
			g.Printf(`final MethodChannel channel = new MethodChannel(registrar.messenger(), "%s");`, )
			g.Printf("channel.setMethodCallHandler(new %s());", )
			g.Outdent()
			g.Printf("}\n")

			g.Printf("@Override")
			g.Printf("public void onMethodCall(MethodCall call, Result result) {\n")
			g.Indent()

			for _, cons := range cons {
				if !g.isConsSigSupported(cons.Type()) {
					g.Printf("// skipped constructor %s.%s with unsupported parameter or return types\n\n", n, cons.Name())
					continue
				}
			}

			for _, f := range fields {
				if t := f.Type(); !g.isSupported(t) {
					g.Printf("// skipped field %s.%s with unsupported type: %s\n\n", n, f.Name(), t)
					continue
				}
				fdoc := doc.Member(f.Name())
				g.javadoc(fdoc)

				g.Indent()
				g.Printf(`if (call.method.equals("get%s")) {\n`)
				g.flutterListFinder()
				g.Printf("result.success(item.get%s());", f.Name())
				g.Printf("return;\n")
				g.Outdent()
				g.Printf("}\n")

				g.Indent()
				g.Printf(`if (call.method.equals("set%s")) {\n`, f.Name())
				g.flutterListFinder()
				g.Printf("item.set%s(call.args[1]));", f.Name())
				g.Printf("result.success();")
				g.Printf("return;\n")
				g.Outdent()
				g.Printf("}\n")

			}

			for _, m := range methods {
				if !g.isSigSupported(m.Type()) {
					g.Printf("// skipped method %s.%s with unsupported parameter or return types\n\n", n, m.Name())
					continue
				}
				g.javadoc(doc.Member(m.Name()))
			}

			g.Indent()
			g.Printf(`if (call.method.equals("destroy")) {\n`)
			g.Printf("list.remove(call.args[0]);")
			g.Printf("result.success();")
			g.Printf("return;\n")
			g.Outdent()
			g.Printf("}\n")

			g.Printf("return result.notImplemented();")

			g.Outdent()
			g.Printf("}\n")

		} else if g.isRn {

			for _, cons := range cons {
				if !g.isConsSigSupported(cons.Type()) {
					g.Printf("// skipped constructor %s.%s with unsupported parameter or return types\n\n", n, cons.Name())
					continue
				}
			}

			for _, f := range fields {
				if t := f.Type(); !g.isSupported(t) {
					g.Printf("// skipped field %s.%s with unsupported type: %s\n\n", n, f.Name(), t)
					continue
				}

				fdoc := doc.Member(f.Name())

				g.javadoc(fdoc)
				g.methodAnnotations()
				g.Printf("public void get%s(String id, Callback successCallback, Callback errorCallback) {\n", f.Name())
				g.Indent()
				g.rnListFinder(false)
				g.Printf("successCallback.invoke(item.get%s());", f.Name())
				g.Outdent()
				g.Printf("}")

				g.javadoc(fdoc)
				g.methodAnnotations()
				g.Printf("public void get%s(String id, Promise promise) {\n", f.Name())
				g.Indent()
				g.rnListFinder(true)
				g.Printf("promise.resolve(item.get%s());", f.Name())
				g.Outdent()
				g.Printf("}")

				g.javadoc(fdoc)
				g.methodAnnotations()
				g.Printf("public void set%s(String id, %s v) {\n", f.Name(), g.javaType(f.Type()))
				g.Indent()
				g.Printf("item.set%s(v);", f.Name())
				g.Outdent()
				g.Printf("}")


			}

			for _, m := range methods {
				if !g.isSigSupported(m.Type()) {
					g.Printf("// skipped method %s.%s with unsupported parameter or return types\n\n", n, m.Name())
					continue
				}
				g.javadoc(doc.Member(m.Name()))
			}
		}


		g.methodAnnotations()
		g.Printf("public void destroy(String id) {\n")
		g.Indent()
		g.Printf("this.list.remove(id);")
		g.Outdent()
		g.Printf("}")
		g.Outdent()
		g.Printf("}")

		g.Printf("*/")


	}


}

func (g *JavaGen) genRCTFuncSignature(o *types.Func, jm *java.Func, hasThis bool) {
	if g.isRn {
		g.Printf("/*")
		success := "successCallback.invoke(this.%s(%s));"
		error := "errorCallback.invoke(e.getMessage());"
		finalArgs := "Callback successCallback, Callback errorCallback"
		if g.usePromises {
			success = "promise.resolve(this.%s(%s));"
			error = "promise.reject(e.getMessage());"
			finalArgs = "Promise promise"
		}
		g.Printf("public final void %s", finalArgs)
		g.Printf("(")
		g.Printf(") {\n")
		g.Indent()
		g.Printf("try {")
		g.Indent()
		g.Printf(success)
		g.Outdent()
		g.Printf("} catch (e Exception) {")
		g.Indent()
		g.Printf(error)
		g.Outdent()
		g.Printf("}")
		g.Outdent()
		g.Printf("}")
		g.Printf("*/")

	}

}

func (g *JavaGen) genRCTVar(o *types.Var) {
	if g.isRn {
		g.Printf("/*")

		// rn callback-style getter
		g.javadoc(doc)
		g.methodAnnotations()
		g.Printf("public static void get%s(Callback successCallback) {\n", o.Name())
		g.Indent()
		g.Printf("successCallback.invoke(this.get%s());", o.Name())
		g.Outdent()
		g.Printf("}\n")

		// rn promise-style getter
		g.javadoc(doc)
		g.methodAnnotations()
		g.Printf("public static void get%s(Promise promise) {\n", o.Name())
		g.Indent()
		g.Printf("promise.resolve(this.get%s());", o.Name())
		g.Outdent()
		g.Printf("}\n")
		g.Printf("*/")
	}

}

func (g *JavaGen) methodAnnotations() {
	if g.isRn {
		g.Printf("//@ReactMethod\n")
	}
}

func (g *JavaGen) methodAsync(f *types.Var, fdoc string) {

	g.Printf("/*")

	if g.isFlutter {

		g.javadoc(fdoc)

	} else if g.isRn {

		g.javadoc(fdoc)
		g.methodAnnotations()

		if g.usePromises {
			g.Printf("public final void get%s(Promise promise) {\n", f.Name())
			g.Indent()
			g.Printf("promise.resolve(this.get%s());\n", f.Name())
			g.Outdent()
			g.Printf("}\n")
		} else {
			g.Printf("public final void get%s(Callback successCallback, Callback errorCallback) {\n", f.Name())
			g.Indent()
			g.Printf("successCallback.invoke(this.get%s());\n", f.Name())
			g.Outdent()
			g.Printf("}\n")
		}

		g.Printf("*/")

	}
}

func (g *JavaGen) listFinder(arg, error string) {
	g.Printf("var item = this.list.get(%s);", arg)
	g.Printf("if (item == null) {")
	g.Indent()
	g.Printf(error)
	g.Printf("return;")
	g.Outdent()
	g.Printf("}")
}

func (g *JavaGen) rnListFinder(promise bool) {
	if promise {
		g.listFinder("id", `promise.reject("not found");`)
	} else {
		g.listFinder("id", `errorCallback.invoke("not found");`)
	}
}

const (
	rnPreamble = `/*import com.facebook.react.ReactPackage;
import com.facebook.react.bridge.NativeModule;
import com.facebook.react.bridge.ReactApplicationContext;
import com.facebook.react.uimanager.ViewManager;
import com.facebook.react.bridge.Callback;
import com.facebook.react.bridge.Promise;
*/
`

)