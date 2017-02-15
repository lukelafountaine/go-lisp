(define <= (lambda (x y)
    (or
        (< x y)
        (= x y))))

(define >= (lambda (x y)
    (or
        (> x y)
        (= x y))))

(define abs (lambda x
        (if (< x 0)
            (* -1 x)
            x)))

(define comparator (lambda fn
    (lambda lst
        (if lst
            (begin
                (define helper (lambda (big lst)
                    (if lst
                        (if (fn (car lst) big)
                            (helper (car lst) (cdr lst))
                            (helper big (cdr lst)))
                        big)))
                (helper (car lst) (cdr lst)))
            (quote ())))))

(define max (comparator >))

(define min (comparator <))

(define map (lambda (fn lst)
    (if lst
        (cons (fn (car lst)) (map fn (cdr lst)))
        (quote ()))))

(define filter (lambda (fn lst)
    (if lst
        (if (fn (car lst))
            (cons (car lst) (filter fn (cdr lst)))
            (filter fn (cdr lst)))
        (quote ()))))

(define reduce (lambda (fn lst val)
    (if lst
        (reduce fn (cdr lst) (fn val (car lst)))
        val)))

(define range (lambda (a b)
    (if (= a b)
        (quote ())
        (cons a (range (+ a 1) b)))))

(define fib (lambda (n)
    (if (< n 2)
        1
        (+ (fib (- n 1)) (fib (- n 2))))))

(define fact (lambda (x)
    (if (= x 0)
        1
        (* x (fact (- x 1))))))

(define count (lambda (item L)
    (if L
        (+ (if (= item (car L)) 1 0) (count item (cdr L)))
        0)))

(define gcd (lambda (a b)
    (if (|| (< a 0) (< b 0))
        ((gcd (abs a) (abs b)))
        (if (< a b)
            (gcd b a)
            (if (= b 0)
                a
                (gcd b (% a b)))))))