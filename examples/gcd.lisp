(define gcd (lambda (a b)
    (if (|| (< a 0) (< b 0))
        ((gcd (abs a) (abs b)))
        (if (< a b)
            (gcd b a)
            (if (== b 0)
                a
                (gcd b (% a b)))))))

(&&
    (== (gcd 13 13) 13)
    (== (gcd 37 600) 1)
    (== (gcd 20 100) 20)
    (== (gcd 624129 2061517) 18913))