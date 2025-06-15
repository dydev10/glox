# TODO

### Lexer
[-] Multiline comment  
[-] Nesting multiline comment  

### AST
[x] Refactor generic classes  

### Parser
[-] Add block statements -> comma operator support  
[-] Add ternary operator  
[-] Report error when binary operator without left operand  

### Interpreter
[-] Support comparison b/w strings and number, and b/w two string   
[-] Implicitly convert other operand to string for concatenation if one operand is string   
[-] Handle divide by zero, give runtime error or define behavior like assigning infinity    

### REPL
[-] Fix REPL to maintain state  
[-] Fix REPL to evaluate expression directly without statement  

### Statements
[-] Throw runtime error when accessing uninitialized vars instead of implicit nil 

### Control Flow
[-] Add break and continue keywords support to early exit loops and if-else blocks  



