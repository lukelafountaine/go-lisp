(define make-account
  (lambda (balance)
    (lambda (amt)
      (begin (set! balance (+ balance amt))
             balance))))

(define account1 (make-account 100.00))
(account1 -20.00)
