(define first car)
(define rest cdr)
(define equal? =)

(define count (lambda (item L)
    (if L
        (+ (if (equal? item (first L)) 1 0) (count item (rest L)))
        0)))

(&&
    (= (count 0 (list 0 1 2 3 0 0) 3))
    (= (count (quote the) (quote (the more the merrier the bigger the better))) 4))

