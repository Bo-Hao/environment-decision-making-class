
from math import *
import random
import numpy as np 
 
#机器人四个参照物
landmarks  = [[20.0, 20.0], [80.0, 80.0], [20.0, 80.0], [80.0, 20.0]]
#地图大小
world_size = 100.0
class robot:
    def __init__(self):
        self.x = random.random() * world_size
        self.y = random.random() * world_size
        self.orientation = random.random() * 2.0 * pi
        #给机器人初试化一个坐标和方向
        self.forward_noise = 0.0;
        self.turn_noise    = 0.0;
        self.sense_noise   = 0.0;
    
    def set(self, new_x, new_y, new_orientation):
		#设定机器人的坐标　方向
        self.x = float(new_x)
        self.y = float(new_y)
        self.orientation = float(new_orientation)
    
    
    def set_noise(self, new_f_noise, new_t_noise, new_s_noise):
        # makes it possible to change the noise parameters
        # this is often useful in particle filters
        #设定一下机器人的噪声
        self.forward_noise = float(new_f_noise);
        self.turn_noise    = float(new_t_noise);
        self.sense_noise   = float(new_s_noise);
    
    
    def sense(self):
		#测量机器人到四个参照物的距离　可以添加一些高斯噪声
        Z = []
        for i in range(len(landmarks)):
            dist = sqrt((self.x - landmarks[i][0]) ** 2 + (self.y - landmarks[i][1]) ** 2)
            dist += random.gauss(0.0, self.sense_noise)
            Z.append(dist)
        return Z
    
    
    def move(self, turn, forward):
        #机器人转向　前进　并返回更新后的机器人新的坐标和噪声大小
        # turn, and add randomness to the turning command
        orientation = self.orientation + float(turn) + random.gauss(0.0, self.turn_noise)
        orientation %= 2 * pi
        
        # move, and add randomness to the motion command
        dist = float(forward) + random.gauss(0.0, self.forward_noise)
        x = self.x + (cos(orientation) * dist)
        y = self.y + (sin(orientation) * dist)
        x %= world_size    # cyclic truncate
        y %= world_size
        
        # set particle
        res = robot()
        res.set(x, y, orientation)
        res.set_noise(self.forward_noise, self.turn_noise, self.sense_noise)
        return res
    
    def Gaussian(self, mu, sigma, x):
        
        # calculates the probability of x for 1-dim Gaussian with mean mu and var. sigma
        return exp(- ((mu - x) ** 2) / (sigma ** 2) / 2.0) / sqrt(2.0 * pi * (sigma ** 2))
    
    
    def measurement_prob(self, measurement):
        
        # calculates how likely a measurement should be
        #计算出的距离相对于正确正确的概率　离得近肯定大　离得远就小
        prob = 1.0;
        for i in range(len(landmarks)):
            dist = sqrt((self.x - landmarks[i][0]) ** 2 + (self.y - landmarks[i][1]) ** 2)
            prob *= self.Gaussian(dist, self.sense_noise, measurement[i])
        return prob

def plot_particle(p):
    import matplotlib.pyplot as plt
    X, Y = [], []
    for i in range(len(p)):
        X.append(p[i].x)
        Y.append(p[i].y)
    
    plt.scatter(X, Y)
    l = np.array(landmarks).T
    plt.scatter(l[0], l[1], c = "red")
    plt.show()


if __name__ == "__main__":
    N = 1000
    #初始化一千个粒子
    p = []
    for i in range(N):
        x = robot()
        x.set_noise(0.05, 0.05, 5.0)
        p.append(x)

    # try 100 steps 
    for r in range(200):
        p2 = []
        for i in range(N):
            p2.append(p[i].move(0.5, 5.0))
        p = p2
        #计算各个粒子的权重
        w = []
        for i in range(N):
            w.append(p[i].measurement_prob(p[i].sense()))


        p3 = []
        #初始化一个数组下标
        index = int(random.random() * N)
        beta = 0.0
        #获取W 中最大的那个数据
        maxw = max(w)
        #运行Ｎ次
        for i in range(N):
            beta += random.random() * 2.0 * maxw
            #beta 变量变大
            #如果beta 小于w[index] 则直接跳过while 循环
            while beta > w[index]:
                #如果beta 大于w[index] 则减去该值　同时index值加1 ％N 就可以循环　知道beta 小于w[index]
                beta -= w[index]
                index = (index + 1) % N
            p3.append(p[index])
        p = p3

    print(len(p))
    plot_particle(p)
    # conclusion: It is not much different between golang and python.
