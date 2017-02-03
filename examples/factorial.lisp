(define fact (lambda (x)
    (if (= x 0)
        1
        (* x (fact (- x 1))))))

(&&
    (= (fact 5) 120)
    (= (fact 10) 3628800))
