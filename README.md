# Doppelkastenschluessel Hex
Doppelkastenschluessel Hex - A slightly simplified version, in hex notation, with support for multiple keys, of the German WWII pencil and paper cipher.

Example usage:  

 $ echo -n 'Hello World.' | dkh > message.txt  
 $ cat message.txt  
48656C6C6F20576F  
726C642E4B4E4C50  

$ dks < message.txt > seriation.txt  
$ cat seriation.txt  
47 82 66 5C 66 C4 62 CE 64 FB 24 0E 54 7C 65 F0  

$ dkt+ -p test -s test > keys.txt  
$ cat keys.txt  
Tagesschlüssel für den 15.12.2024 (UTC)  

Tagesschlüssel 1:  
Kasten: A    Kasten: B  
7 D 6 4      8 4 2 3  
2 A 8 C      E 9 A 6  
0 1 3 9      F B 7 1  
5 E B F      D 5 C 0   

$ dkk -s 1 < seriation.txt > ciphertext.txt  
$ cat ciphertext.txt  
33 A6 2C DB 2C 6D 26 62 2D 01 ED F2 DD 8B 2E 0F  
