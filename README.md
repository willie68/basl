# Basl Basic Seriell Language

Diese Sprache ist inspiriert von SIMPL, VTL-02, txtzyme und anderen sehr kleinen Sprachen. Das Hauptaugenmerk der Sprache ist einerseits die einfache Lernbarkeit. Andererseits sollen damit die üblichen Aufgaben im Umfeld von Microcontrollern im Modellbau adressiert werden. Und ähnlich wie die TPS soll die Programmierung auch im eingebauten Zustand ohne Zuhilfenahme eines PC ermöglicht werden. 

## Sprachgrundlagen

Die Sprache ist eine interpretierende REPL Sprache. (REPL: read-evaluate-print-loop) D.h. der Interpreter liest zunächst von der seriellen Konsole eine Zeile, interpretiert diese, und gibt dann das Ergebnis wieder aus. Danach fängt der Zyklus von vorne an. 

Der Interpreter ist eine sog. Stackmaschine, d.h. alle Befehle arbeiten mit dem Stack, manipulieren diesen, entnehmen Werte oder legen neue Werte darauf. Ein zusätzlicher Speicher kann angesprochen werden. Größe des Stacks und des Speichers sind implementierungsabhängig. Auch die Zuordnung der Pins zu analogen oder digitalen Ein/Ausgaben finden sich in der jeweiligen Implementierung. 

Die grundlegenden Befehle der Sprache sind als Kleinbuchstaben definiert. Für eigene "Routinen" oder "Funktionen" stehen damit die Kleinbuchstaben zur Verfügung. Die Sprache verwendet einen Parameterstack, um Parameter an Befehle, Routinen oder Funktionen zu übergeben. 

Beispiel: 

Die Zeile 

`100 200 + [cr]`

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

 **\>** und **<** verschieben den aktuellen  

### Blöcke, Bedingungen, Schleifen 

mit **{}** können Blöcke erzeugt werden. Diese gelten dann bei Bedingungen und schleifen als ein Befehl 

Mit **?** wird ein Wert vom Stack genommen und damit ein Test ausgeführt. Ist der aktuelle Wert auf dem Stack > 0 wird der erste nachfolgende Befehl/Block ausgeführt. Ist der Wert = 0, wird der nächste Befehl/Block übersprungen und der übernächste ausgeführt.

Mit **=** werden 2 Werte vom Stack genommen und auf Gleichheit getestet. Ist das Ergebniss wird der erste nachfolgende Befehl/Block ausgeführt. Ist das Ergebnis falsch, wird der nächste Befehl/Block übersprungen und der übernächste ausgeführt.  

Ähnliches gilt für **>** und **<**. 

Schleifen werden mit **#** gestartet. Der 1. Parameter ist die max. Anzahl der Durchläufe. Auf den Zähler der Schleife kann mittels **k** zugegriffen werden 

Beispiel: `20 # { k p}` gibt die Werte `1 2 3 4 5 .. 20` aus. 

**b** ist ein Break, d.h. die Schleife wird direkt verlassen. ist gerade keine Schleife aktiv, wird auf die Shell gesprungen.

**c** bedeutet Continue, d.h. die Schleife wird mit dem nächsten Index weiter gemacht. Alle Befehle nach c werden ignoriert.

### Input, Output

**o** gibt einen Wert auf einen Ausgabepin aus. `13 1 o [cr]` auf Pin 13 wird eine 1 ausgegeben. Je nach Konfiguration des Pins kann das auch ein analoger (PWM) Wert sein. 

**i** legt den Wert eines Eingabepins auf den Stack. `13 i [cr]` Das aktuelle Zustand von Pin13 wird auf den Stack gelegt. Je nach Konfiguration kann der Wert auch ein analoger Wert (ADC) sein. 

**j** misst die Impulslänge an einem Pin. 

**p** gibt des aktuellen Wert auf der Schnittstelle aus.

**_** gibt den nachfolgenden Text bis zum CR aus. `_ Hallo RC Simple [cr]`  gibt "Hallo RC Simple" aus.

### Kommandos

**d** Delay, also eine Wartezeit, der Parameter gibt die Anzahl der ms an.

**t** Tone, erzeugt einen Ton, der Parameter gibt die Frequenz an. 0 bedeutet ausschalten. Beispiel: 

`440 t 1000 d 0 t` erzeugt einen 1 sekündigen 440 Hz Ton.

### Eigene Routinen und Befehle

Ein eigener Befehl wir mit eine **:** eingeleitet. Danach folgt ein Großbuchstabe, der den Namen der Routine vorgibt. danach folgen die Befehle. Die Definition endet mit einem **;** Nur diese Routinen werden dauerhaft gespeichert.  

Eine Besonderheit ist die Routine mit A. Diese wird nach dem Start des Systemes automatisch ausgeführt.

### Befehlsübersicht

| Zeichen | Bedeutung                                                    | Parameter                         | Zeichen | Bedeutung                                                    | Parameter                                 |
| ------- | ------------------------------------------------------------ | --------------------------------- | ------- | ------------------------------------------------------------ | ----------------------------------------- |
| a       |                                                              |                                   | b       | Break, ein Block wird abgebrochen,<br />ist gerade kein Block aktiv, wird auf die Shell gesprungen |                                           |
| c       | Continue, in einem Schleifenblock wird mit dem nächsten Indexwert weiter gemacht. |                                   | d       | delay                                                        | 1. Anzahl der ms                          |
| e       |                                                              |                                   | f       |                                                              |                                           |
| g       |                                                              |                                   | h       | help, zeigt alle Befehle an                                  |                                           |
| i       | input from pin                                               | 1. Pinnummer                      | j       | Pulse in, misst die Pulsweite am Pin in ms                   | 1. Pin                                    |
| k       | aktueller Wert in einer Schleife                             |                                   | l       |                                                              |                                           |
| m       |                                                              |                                   | n       |                                                              |                                           |
| o       | output to pin                                                | 1. Pinnumer<br />2. Wert          | p       | gibt den aktuellen Wert auf der Schnittstelle aus            | 1. Wert                                   |
| q       | gibt alle Unterprogramme aus                                 |                                   | r       | Restore, holt einen Wert aus einer Speicherzelle             | 1. Speicherzelle                          |
| s       | Save, speichert einen Wert auf eine Speicherstelle           | 1. Speicherzelle<br />2. Wert     | t       | tone                                                         | 1.  Frequenz 0=Off                        |
| u       |                                                              |                                   | v       |                                                              |                                           |
| w       |                                                              |                                   | x       |                                                              |                                           |
| y       |                                                              |                                   | z       |                                                              |                                           |
| !       |                                                              |                                   |         |                                                              |                                           |
| "       | DUP, obersten Stackwert duplizieren                          |                                   | §       |                                                              |                                           |
| /       | Division                                                     | 1. Wert<br />2. Wert              | $       |                                                              |                                           |
| %       | Modulus                                                      | 1. Wert<br />2. Wert              | &       | AND                                                          | 1. Wert<br />2. Wert                      |
| ()      |                                                              |                                   | =       | Skip if equal                                                | 1. Wert<br />2. Wert                      |
| []      |                                                              |                                   | {}      | Block definition                                             |                                           |
| +       | Addition                                                     | 1. Wert<br />2. Wert              | ?       | Skip if 0                                                    | 1. Wenn =0 dann Befehl/Block überspringen |
| *       | Multiplikation                                               | 1. Wert<br />2. Wert              | ~       | NOT                                                          | 1. Wert                                   |
| -       | Subtraktion                                                  | 1. Wert<br />2. Wert              | _       | Gibt einen Text auf der Schnittstelle aus.                   |                                           |
| #       | Start einer Schleife                                         | 1. Anzahl der Schleifendurchgänge | : ;     | Start und Ende einer eigenen Definition                      |                                           |
| .       | print stacksize                                              |                                   | ,       | print stack                                                  |                                           |
| °       |                                                              |                                   | @       |                                                              |                                           |
| €       |                                                              |                                   | ^       | XOR                                                          | 1. Wert<br />2. Wert                      |
| \|      | OR                                                           | 1. Wert<br />2. Wert              | >       | Skip if Greater than                                         | 1. Wert<br />2. Wert                      |
| '       | DROP, obersten Stackwert verwerfen                           |                                   | <       | Skip if lesser than                                          | 1. Wert<br />2. Wert                      |

## Apendix

SIMPL: https://github.com/monsonite/SIMPL

VTL-02: https://altairclone.com/downloads/roms/VTL-2%20(Very%20Tiny%20Language)/VTL-2%20Manual.pdf 

txtzyme: https://github.com/WardCunningham/Txtzyme
