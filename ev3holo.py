#!/usr/bin/env python3
import sys
import ev3dev.ev3 as ev3
from math import cos,sin,pi
from sympy.matrices import Matrix
from sympy import pprint
sys.displayhook = pprint

aA = 0
aB = 120
aC = 240

matrix = Matrix([
    [cos(aA*pi/180),cos(aB*pi/180),cos(aC*pi/180)],
    [sin(aA*pi/180),sin(aB*pi/180),sin(aC*pi/180)],
    [1,1,1]
]).inv()
pprint(matrix)

mA = ev3.Motor('outA')
mB = ev3.Motor('outB')
mC = ev3.Motor('outC')

maxA = mA.max_speed
maxB = mB.max_speed
maxC = mC.max_speed
print(maxA)
print(maxB)
print(maxC)

def holo(m):
    pprint(m)
    f = matrix * m
    pprint(f)

    speedA = maxA * f[0]
    speedB = maxB * f[1]
    speedC = maxC * f[2]
    print(speedA)
    print(speedB)
    print(speedC)

    if mA.wait_until_not_moving() and mB.wait_until_not_moving() and mC.wait_until_not_moving():
        mA.run_timed(time_sp=1000, speed_sp=int(speedA))
        mB.run_timed(time_sp=1000, speed_sp=int(speedB))
        mC.run_timed(time_sp=1000, speed_sp=int(speedC))

holo(Matrix([ [1], [0],[0]]))
holo(Matrix([ [0], [1],[0]]))
holo(Matrix([[-1], [0],[0]]))
holo(Matrix([ [0],[-1],[0]]))
