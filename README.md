# Basl, Basic Serial Language

Diese Sprache ist inspiriert von SIMPL, VTL-02, txtzyme und anderen sehr kleinen Sprachen. Das Hauptaugenmerk der Sprache ist einerseits die einfache Lernbarkeit. Andererseits sollen damit die üblichen Aufgaben im Umfeld von Microcontrollern im Modellbau adressiert werden. Und ähnlich wie die TPS soll die Programmierung auch im eingebauten Zustand ohne Zuhilfenahme eines PC ermöglicht werden. 

## Sprachgrundlagen

Die Sprache ist eine interpretierende REPL Sprache. (REPL: read-evaluate-print-loop) D.h. der Interpreter liest zunächst von der seriellen Konsole eine Zeile, interpretiert diese, und gibt dann das Ergebnis wieder aus. Danach fängt der Zyklus von vorne an. 

Die Eingaben können aber auch direkt aus einem Speicher (EEPROM, Flash) stammen. Der Interpreter selber ist eine sog. Stackmaschine, d.h. alle Befehle arbeiten mit dem Stack, manipulieren diesen, entnehmen Werte oder legen neue Werte darauf. Auch steht ein zusätzlicher Speicher zur Verfügung. Die Größe des Stacks und des Speichers sind implementierungsabhängig. Auch die Zuordnung der Pins zu analogen oder digitalen Ein/Ausgaben finden sich in der jeweiligen Implementierung. 

Die grundlegenden Befehle der Sprache sind als Kleinbuchstaben definiert. Für eigene "Routinen" oder "Funktionen" stehen damit die Großbuchstaben zur Verfügung. Die Sprache verwendet einen Parameterstack, um Parameter an Befehle, Routinen oder Funktionen zu übergeben. 

Beispiel: 

Die Zeile 

`100 200+[cr]`

berechnet 100 + 200 und legt dann das Ergebnis, 300, wieder auf den Stack.



## Befehle 

### Mathematik

**\+** Addiert 2 Werte und legt das Ergebnis wieder auf den Stack. Beispiel: 

`100 200 + [cr]`

berechnet 100 + 200 und legt dann das Ergebnis, 300, wieder auf den Stack.   

**\-** Subtrahiert 2 Werte und legt das Ergebnis wieder auf den Stack. Beispiel:

`200 100 - [cr]`

berechnet 200 - 100 und legt dann das Ergebnis, 100, wieder auf den Stack.   

Ähnlich funktionieren: 

**\*** Multiplikation

**/** Division

**%** Divisionsrest

**&** AND, mathematisches Und

**|** OR, mathematisches Oder

**^** XOR, mathematisches exclusives Oder

**~** NOT, also das nicht, das jedoch nur auf den 1. Stack Parameter angewendet wird.

Mit **=** werden 2 Werte vom Stack genommen und auf Gleichheit getestet. Ist das Ergebnis (1 oder 0) wird dann wieder auf den Stack geparkt.

Ähnliches gilt für **>** und **<**. 

### Bedingungen, Schleifen 

Bedingungen und Schleifen werden über **{}** angelegt. **{** leitet den Block ein. Dabei wird ein Wert vom Stack geholt. Bei 0 wird der Block ignoriert, bei Werten >0 wird der Block genau so oft ausgeführt. 

Auf den Zähler der Schleife kann mittels **k** zugegriffen werden.

Beispiel: `20{ k p}` gibt die Werte `1 2 3 4 5 .. 20` aus. 

**b** ist ein Break, d.h. die Schleife wird direkt verlassen. ist gerade keine Schleife aktiv, wird auf die Shell gesprungen.

**c** bedeutet Continue, d.h. die Schleife wird mit dem nächsten Index weiter gemacht. Alle Befehle nach c werden ignoriert.

Verschachtelte Schleifen sind nicht erlaubt. Das Ergebnis ist undefiniert.

### Input, Output

**o** gibt einen Wert auf einen Ausgabepin aus. `13 1 o [cr]` auf Pin 13 wird eine 1 ausgegeben. Je nach Konfiguration des Pins kann das auch ein analoger (PWM) Wert sein. 

**i** legt den Wert eines Eingabepins auf den Stack. `13 i [cr]` Das aktuelle Zustand von Pin13 wird auf den Stack gelegt. Je nach Konfiguration kann der Wert auch ein analoger Wert (ADC) sein. 

**j** misst die Impulslänge an einem Pin. 

**p** gibt des aktuellen Wert auf der Schnittstelle aus.

**_** gibt den nachfolgenden Text bis zum nächsten **_** aus. `_ Hallo RC Simple_`  gibt "Hallo RC Simple" aus.

**@** können Pins definiert werden. Je nach Implementierung können hier verschiedene Einstellungen vorgenommen werden. In der GoLang implementierung besteht der Config Text aus verschiedenen Buchstaben. Die Position gibt den Pin Index an. 
**i** Digital Input
**o** Digital Output
**a** analog Input
**p** PWM (analog) Output
**s** Servo Output
Beispiel: `@ iiiiooooippixoaa` ist die typische  Arduino TPS Konfiguration

**$** gibt die aktuelle Pin Konfiguration aus

### Kommandos

**d** Delay, also eine Wartezeit, der Parameter gibt die Anzahl der ms an.

**t** Tone, erzeugt einen Ton, der Parameter gibt die Frequenz an. 0 bedeutet ausschalten. Beispiel: 

`440 t 1000 d 0 t` erzeugt einen 1 sekündigen 440 Hz Ton.

### Stack

**[xxxx]** legt die Nummer als Wert auf den Stack, z.B. `100` legt die `100` auf den Stack, `12 23 34 45` ergibt einen Stack mit `12 23 34 45`

**p** holt den obersten Wert des Stacks und gibt ihn auf der Console aus

**"** dupliziert den obersten Stackwert

**'** löscht den obersten Wert des Stacks

**z** löscht den kompletten Stack

**,** gibt alle Stackwerte aus, verändert aber nicht den Stack

**.** gibt die Anzahl der Wert auf den Stack, verändert aber nicht den Stack

**!** vertauscht die beiden ersten Werte auf dem Stack

### Eigene Routinen und Befehle

Ein eigener Befehl wir mit eine **:** eingeleitet. Danach folgt ein Großbuchstabe, der den Namen der Routine vorgibt. danach folgen die Befehle. Die Definition endet mit einem **;** Nur diese Routinen werden dauerhaft gespeichert.  

Eine Besonderheit ist die Routine mit A. Diese wird nach dem Start des Systems automatisch ausgeführt.

Das System ist nicht rekursiv, d.h. Blöcke (Schleifen und Conditions) dürfen nicht verschachtelt sein. 

### Befehlsübersicht

| Zeichen | Bedeutung                                                    | Parameter                     | Zeichen | Bedeutung                                                    | Parameter                         |
| ------- | ------------------------------------------------------------ | ----------------------------- | ------- | ------------------------------------------------------------ | --------------------------------- |
| a       |                                                              |                               | b       | Break, ein Block wird abgebrochen,<br />ist gerade kein Block aktiv, wird auf die Shell gesprungen |                                   |
| c       | Continue, in einem Schleifenblock wird mit dem nächsten Indexwert weiter gemacht. |                               | d       | delay                                                        | 1. Anzahl der ms                  |
| e       |                                                              |                               | f       |                                                              |                                   |
| g       |                                                              |                               | h       | help, zeigt alle Befehle an                                  |                                   |
| i       | input from pin                                               | 1. Pinnummer                  | j       | Pulse in, misst die Pulsweite am Pin in ms                   | 1. Pin                            |
| k       | aktueller Wert in einer Schleife                             |                               | l       |                                                              |                                   |
| m       |                                                              |                               | n       | number: get a number from the console, in interactive, simply type the number after the n and press Enter, e.g. n123 |                                   |
| o       | output to pin                                                | 1. Pinnumer<br />2. Wert      | p       | gibt den aktuellen Wert aus                                  | 1. Wert                           |
| q       | gibt alle Unterprogramme aus                                 |                               | r       | Restore, holt einen Wert aus einer Speicherzelle             | 1. Speicherzelle                  |
| s       | Save, speichert einen Wert auf eine Speicherstelle           | 1. Speicherzelle<br />2. Wert | t       | tone                                                         | 1.  Frequenz 0=Off                |
| u       |                                                              |                               | v       |                                                              |                                   |
| w       |                                                              |                               | x       |                                                              |                                   |
| y       |                                                              |                               | z       | Clear stack                                                  |                                   |
|         |                                                              |                               |         |                                                              |                                   |
| "       | DUP, obersten Stackwert duplizieren                          |                               | !       | SWAP, vertauscht die beiden oberen Stackwerte                |                                   |
| /       | Division                                                     | 1. Wert<br />2. Wert          | $       | output pin configuration                                     |                                   |
| %       | Modulus                                                      | 1. Wert<br />2. Wert          | &       | AND                                                          | 1. Wert<br />2. Wert              |
| ()      |                                                              |                               | =       | Equal                                                        | 1. Wert<br />2. Wert              |
| []      |                                                              |                               | {}      | Schleife, Bedingung                                          | 1. Anzahl der Schleifendurchgänge |
| +       | Addition                                                     | 1. Wert<br />2. Wert          | ?       |                                                              |                                   |
| *       | Multiplikation                                               | 1. Wert<br />2. Wert          | ~       | NOT                                                          | 1. Wert                           |
| -       | Subtraktion                                                  | 1. Wert<br />2. Wert          | _       | Gibt einen Text auf der Schnittstelle aus, bis zum nächsten _. |                                   |
| #       |                                                              |                               | : ;     | Start und Ende einer eigenen Definition                      |                                   |
| .       | print stacksize                                              |                               | ,       | print stack                                                  |                                   |
| ^       | XOR                                                          | 1. Wert<br />2. Wert          | @       | Config: hier kann die aktuelle Konfigurtion abgelegt werden. Gilt bis zum nächsten CR |                                   |
| \|      | OR                                                           | 1. Wert<br />2. Wert          | >       | Greater than                                                 | 1. Wert<br />2. Wert              |
| '       | DROP, obersten Stackwert verwerfen                           |                               | <       | Lesser than                                                  | 1. Wert<br />2. Wert              |

## Apendix

SIMPL: https://github.com/monsonite/SIMPL

VTL-02: https://altairclone.com/downloads/roms/VTL-2%20(Very%20Tiny%20Language)/VTL-2%20Manual.pdf 

txtzyme: https://github.com/WardCunningham/Txtzyme
