from pygamehelper import *
from pygame import *
from vec2d import vec2d
from math import e, pi, cos, sin, sqrt
from random import uniform

# razmerite v tazi simulaciq sa 1:50

class Drone:
    '''Class for workers!'''
    pos = vec2d(0, 0)
    target = []
    single_t = True
    speed = 5
        
    def update(self):
        dir = self.target[0] - self.pos      
        if dir.length >= self.speed:
            dir.length = self.speed*((dir.length//50)+1)
            self.pos = vec2d(int(self.pos[0] + dir[0]), int(self.pos[1] + dir[1]))
        else:
            if len(self.target) > 1:
                temptargets = self.target[1:]
                self.target = temptargets
        print(self.target)

    def colider (self, odrone, targ):
        dist = self.pos.get_distance(odrone.pos)
        if dist < 42:
            overlap = 42 - dist
            ndir = odrone.pos - self.pos
            ndir.length = overlap
            if self == targ.selected:
                self.pos = self.pos - ndir
            elif odrone == targ.selected:
                odrone.pos = odrone.pos + ndir
            else:
                ndir.length =  ndir.length / 2
                odrone.pos = odrone.pos + ndir
                self.pos = self.pos - ndir

    def movedrone(self, pos):
        targets = []
        dtarget = vec2d(pos)
        targets.append(dtarget)
        self.target = targets


class Starter(PygameHelper):
    drones = []
    others = []
    friendscount = 3
    friend = Drone()

    def __init__(self, size=(1300, 700), fill=((255, 255, 255))):    
        self.w, self.h = size
        super().__init__(size=(self.w, self.h), fill=((255,255,255)))

        targetpos = [0, 0]
        for drone in self.drones:
            targetpos[0] += drone.pos[0]
            targetpos[1] += drone.pos[1]
        targetpos[0] /= self.friendscount
        targetpos[1] /= self.friendscount
        self.endpoint = vec2d(targetpos[0], targetpos[1])

        self.friend.pos = vec2d(int(uniform(0, self.w - 20)), int(uniform(0, self.h - 20)))
        target = vec2d( self.friend.pos)
        self.friend.target.append(target)
        self.selected = self.friend

        for i in range(self.friendscount): 
            tempagent = Drone()
            tempagent.pos = vec2d(int(uniform(0, self.w - 20)), int(uniform(0, self.h - 20)))
            target = vec2d(tempagent.pos)
            tempagent.target.append(target)
            self.drones.append(tempagent)
    
    def reset(self):
        self.others = []

        for drone in self.drones:
            drone.pos = vec2d(int(uniform(0, self.w - 20)), int(uniform(0, self.h - 20)))
            target = vec2d(drone.pos)
            drone.target.append(target)

        self.friend.pos = vec2d(int(uniform(0, self.w - 20)), int(uniform(0, self.h - 20)))
        target = vec2d(self.friend.pos)
        self.friend.target = []
        self.friend.target.append(target)

        targetpos = [0, 0]
        for drone in self.drones:
            targetpos[0] += drone.pos[0]
            targetpos[1] += drone.pos[1]
        targetpos[0] /= self.friendscount
        targetpos[1] /= self.friendscount
        self.endpoint = vec2d(targetpos[0], targetpos[1])

    def update(self):
        self.selected.update()
        for drone in self.drones + self.others:
            if not self.selected == drone:
                self.selected.colider(drone, self)
        
    def keyUp(self, key):
        if key == 100:
            self.reset()
        
    def mouseUp(self, button, pos):
        if button == 3:
            self.selected.movedrone(self.endpoint)
        elif button == 1:
            agent = Drone()
            agent.pos = vec2d(pos)
            target = vec2d(agent.pos)
            agent.target.append(target)
            self.others.append(agent)

        
    def mouseMotion(self, buttons, pos, rel):
       if buttons[0]:
            agent = Drone()
            agent.pos = vec2d(pos)
            target = vec2d(agent.pos)
            agent.target.append(target)
            self.others.append(agent)
        
    def draw(self):
        self.screen.fill((0, 0, 0))

        pygame.draw.circle(self.screen, (0, 150, 0),(int(self.selected.pos[0]),int(self.selected.pos[1])), 20, 3)
        pygame.draw.circle(self.screen, (155, 155, 155),(int(self.endpoint[0]),int(self.endpoint[1])), 2)

        for drone in self.drones:
            pygame.draw.circle(self.screen, (150, 0, 0),(int(drone.pos[0]),int(drone.pos[1])), 20, 3)
        for other in self.others:
            pygame.draw.circle(self.screen, (0, 0, 150),(int(other.pos[0]),int(other.pos[1])), 20, 3)
        
s = Starter()
s.mainLoop(40)
