(dap-register-debug-template
  "Go: Debug Boilerplate"
  (list :type "go"
        :cwd (lsp-workspace-root)
        :request "launch"
        :name "Go: Debug boilerplate"
        :mode "auto"
        :program "cmd/service/main.go"
        :buildFlags nil
        :args nil
        :env nil))
