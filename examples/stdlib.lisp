(define map (lambda (fn lst)
    (if lst
        (cons (fn (car lst)) (map fn (cdr lst)))
        (quote ()))))

(define range (lambda (a b)
    (if (= a b)
        (quote ())
        (cons a (range (+ a 1) b)))))

(define fact (lambda (x)
    (if (= x 0)
        1
        (* x (fact (- x 1))))))

(define count (lambda (item L)
    (if L
        (+ (if (= item (car L)) 1 0) (count item (cdr L)))
        0)))