def electron_app_(ctx):
    ctx.actions.run(
      executable = ctx.executable._electron_bundle_tool,
      inputs = [ctx.file.main_js, ctx.file.index_html, ctx.file._electron_release],
      arguments = [
        ctx.outputs.apptar.path,
        ctx.attr.app_name,
        ctx.file.main_js.path,
        ctx.file.index_html.path,
        ctx.file._electron_release.path,
      ],
      outputs = [ctx.outputs.apptar],
    )

    ctx.actions.expand_template(
      template=ctx.file._electron_app_script_tpl,
      output=ctx.outputs.run,
      substitutions = {
            "{{app}}": ctx.outputs.apptar.short_path,
            "{{name}}": ctx.attr.app_name,
        },
      is_executable=True,
    )

    return DefaultInfo(
        executable = ctx.outputs.run,
        files = depset([ctx.outputs.apptar]),
        runfiles = ctx.runfiles(files = [ctx.outputs.apptar]),
    )

electron_app = rule(
    implementation = electron_app_,
    executable = True,
    attrs = {
      "app_name": attr.string(),
      "main_js": attr.label(allow_single_file = True),
      "index_html": attr.label(allow_single_file = True),
      "_electron_bundle_tool": attr.label(
          executable = True,
          cfg = "host",
          allow_files = True,
          default = Label("//:bundle"),
      ),
      "_electron_app_script_tpl": attr.label(
          allow_single_file = True,
          default = Label("//:run.sh.tpl"),
      ),
      "_electron_release": attr.label(
          allow_single_file = True,
          default = Label("@electron_release//file"),
      ),
    },
    outputs = {
        "apptar": "%{name}.tar",
        "run": "%{name}.sh",
    },
)
