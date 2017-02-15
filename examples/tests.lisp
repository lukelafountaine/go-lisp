(and
    (= (count 0 (list 0 1 2 3 0 0)) 3)
    (= (count (quote the) (quote (the more the merrier the bigger the better))) 4)
    (= (gcd 13 13) 13)
    (= (gcd 37 600) 1)
    (= (gcd 20 100) 20)
    (= (gcd 624129 2061517) 18913)
    (= (map fib (range 0 10)) (list 1 1 2 3 5 8 13 21 34 55)))
