class Foo:
    def __init__(self, x) -> None:
        self.x = x

def inner1(f: Foo):
    f.x = 20

def inner2(f: Foo):
    f = Foo(30)

def outer():
    f = Foo(10)
    inner1(f)
    print(f.x) # 20
    inner2(f)
    print(f.x) # 20
    g = None
    inner2(g)
    print(g is None)
outer()