#include <avr/pgmspace.h>

#define prt(S) \
  Serial.print(S);
#define prtln(S) \
  Serial.println(S);
#define ln() \
  Serial.println();


char c;
int nu = 0;
bool inNumber;
#define TONEPIN 10
const int MS = 128;
uint16_t mem[MS]; // memory for 128 numbers
const int SZ = 32;
uint16_t stack[SZ] = {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} ; // stack for 32 numbers
int8_t sp = 0;
uint16_t v  = 0;
uint16_t a = 0;
unsigned long vl;
byte pins[20] = {'i', 'i', 'i', 'i', 'o', 'o', 'o', 'o', 'i', 'p', 'p', 'i', 'x', 'o', 'a', 'a', 'a', 'a', 'a', 'a'};

void setup() {
  Serial.begin(115200);
  prt(F("Basl V0.1"));
  ln();
  pinMode(13, 1);
  inNumber = false;
  nu = 0;
  for (byte x = 0; x < 4; x++) {
    pinMode(x, INPUT_PULLUP);
    pinMode(x + 4, OUTPUT);
  }
  pinMode(8, INPUT_PULLUP);
  pinMode(9, OUTPUT);
  pinMode(10, OUTPUT);
  pinMode(11, INPUT_PULLUP);
  pinMode(13, OUTPUT);
}

bool first;
void loop() {
  first = true;
  while (Serial.available() == 0) {
    if (first)  {
      prt(":");
      first = false;
    }
    delay(10);
  }
  while (Serial.available() > 0) {
    c = Serial.read();
    if ((c >= '0') && (c <= '9')) {
      nu = nu * 10 + (c - '0');
      inNumber = true;
    } else {
      if (inNumber) {
        push(nu);
        inNumber = false;
        nu = 0;
      }
    }
    if ((c == ' ') || (c == '\t') || (c == '\r') || (c == '\n')) {
      break;
    }
    prt(c);
    switch (c) {
      case '0':
      case '1':
      case '2':
      case '3':
      case '4':
      case '5':
      case '6':
      case '7':
      case '8':
      case '9':
        break;
      case '.':
        // print stack size
        ln(); prt(sp);
        break;
      case ',':
        // Print stack
        ln(); prt(F("["));
        if (sp > 0) {
          for (int8_t x = sp; x > 0; x--) {
            prt(stack[x - 1]);
            if (x > 1) prt(F(", "));
          }
        }
        prt("]");
        break;
      case '\"':
        // dupe first stack element
        v = pop();
        push(v);
        push(v);
        break;
      case '\'':
        // drop stack element
        v = pop();
        break;
      case '!':
        // swap first two stack element
        v = pop();
        a = pop();
        push(v);
        push(a);
        break;
      case '_':
        ln();
        while ((c = Serial.read()) != '_') {
          Serial.print(c);
        }
        ln();
        break;
      case '&':
      case '|':
      case '^':
      case '+':
      case '-':
      case '*':
      case '/':
      case '%':
        math(c);
        break;
      case 'd':
        // delay command
        v = pop();
        ln(); prt(F("delay ")); prt(v); prtln(F("ms"));
        delay(v);
        break;
      case 'h':
        // show help
        showHelp();
        break;
      case 'i':
        // delay command
        a = pop();
        if ((a >= 0) && (a < 20) && (pins[a] == 'i')) {
          v = digitalRead(a);
          push(v);
        } else if ((a >= 14) && (a < 20) && (pins[a] == 'a')) {
          v = analogRead(a);
          push(v);
        } else {
          push(0);
        }
        break;
      case 'j':
        // delay command
        v = pop();
        if ((a >= 0) && (a < 20) && (pins[a] == 'i')) {
          vl = pulseIn(v, HIGH, 1000l * 1000); // timeout max 10sec.
          v = vl / 1000;
        }
        push(v);
        break;
      case 'o':
        // delay command
        a = pop();
        v = pop();
        if ((a >= 0) && (a <= 20) && (pins[a] == 'o')) {
          digitalWrite(a, v > 0);
        }
        if ((a >= 0) && (a < 14) && (pins[a] == 'p')) {
          analogWrite(a, v);
        }
        break;
      case 'p':
        // get stack value
        v = pop();
        ln(); prt("value:"); prtln(v);
        break;
      case 'r':
        a = pop();
        if ((a >= 0) || (a < MS)) {
          push(mem[a]);
        }
        break;
      case 's':
        a = pop();
        v = pop();
        if ((a >= 0) || (a < MS)) {
          mem[a] = v;
        }
        break;
      case 't':
        v = pop();
        if (v > 0) {
          tone(TONEPIN, v);
        } else {
          noTone(TONEPIN);
        }
        break;
      case 'z':
        // clear stack
        sp = 0;
        break;
      case 'a':
      case 'b':
      case 'c':
      case 'e':
      case 'f':
      case 'g':
      case 'k':
      case 'l':
      case 'm':
      case 'n':
      case 'q':
      case 'u':
      case 'v':
      case 'w':
      case 'x':
      case 'y':
      default :
        ni();
        break;
    }
  }
}

void math(char c) {
  v = pop();
  a = pop();
  switch (c) {
    case '&':
      push(a & v);
      break;
    case '|':
      push(a | v);
      break;
    case '^':
      push(a ^ v);
      break;
    case '+':
      push(a + v);
      break;
    case '-':
      push(a - v);
      break;
    case '*':
      push(a * v);
      break;
    case '/':
      push(a / v);
      break;
    case '%':
      push(a % v);
      break;
  }
}

void push(int v) {
  if (sp <= SZ ) {
    stack[sp] = v;
    sp++;
  } else sp = SZ;
}

int pop() {
  if (sp > 0) {
    sp--;
    return stack[sp];
  }
  if (sp < 0) sp = 0;
  return 0;
}

void ss() {
  prt("s:"); prt(sp);
  prt("[");
  if (sp > 0) {
    for (int8_t x = sp; x > 0; x--) {
      prt(stack[x - 1]);
      if (x > 1) {
        prt(", ");
      }
    }
  }
  prt("]");
  ln();
}

void ni() {
  ln(); prtln(F("not implemented"));
}

void showHelp() {
  prtln(F("Help"));
  prtln(F("[#]: push # to stack"));
  prtln(F("b: break actual block"));
  prtln(F("c: continue with next interation in loop"));
  prtln(F("d: delay in ms"));
  prtln(F("h: print help"));
  prtln(F("i: input from pin"));
  prtln(F("j: read pulse length"));
  prtln(F("k: push loop index"));
  prtln(F("n: input a number"));
  prtln(F("o: output to pin"));
  prtln(F("@: set pin configuration"));
  prtln(F("$: output pin configuration"));
  prtln(F("r: retrive value from address"));
  prtln(F("s: store value to address"));
  prtln(F("t: tone, 0=off"));
  prtln(F("\": dupe stack value"));
  prtln(F("': drop stack value"));
  prtln(F("§: swap first 2 values on stack"));
  prtln(F("°: clear stack"));
  prtln(F("&: AND"));
  prtln(F("|: OR"));
  prtln(F("^: XOR"));
  prtln(F("~: NOT"));
  prtln(F("+: ADD"));
  prtln(F("-: DEL"));
  prtln(F("*: MUL"));
  prtln(F("/: DIV"));
  prtln(F("%: MOD"));
  prtln(F("p: print value"));
  prtln(F("q: print all user functions"));
  prtln(F(".: print stack size"));
  prtln(F(",: print stack"));
  prtln(F("_: print text till next _"));
  prtln(F("=: skip if not equal"));
  prtln(F("?: skip if not null"));
  prtln(F(">: skip if not greater than"));
  prtln(F("<: skip if not lesser than"));
  prtln(F("{..}}: defining a block"));
}
