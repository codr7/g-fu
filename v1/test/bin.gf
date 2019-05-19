(test (= (int 0xff) 255))
(test (= (inc 0xfe) 0xff))

(let (b (new-bin 1))
  (test (= (# b 0) 0x00))
  (set (# b 0) 0xff)
  (test (= (# b 0) 0xff)))