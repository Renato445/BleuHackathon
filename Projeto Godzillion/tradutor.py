#!/usr/bin/env python
# -*- coding: utf-8 -*-
import os
#import commands

arq = open('saida.txt', 'r')
texto = arq.readlines()

for linha in texto :
    argumento = linha
    if  argumento != "null":
    	os.system("translate-cli -t pt " + argumento + " -o")
    
arq.close()

