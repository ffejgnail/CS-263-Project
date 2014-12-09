Project Cage
============

A virtual environment for artificial animal(animat) to live in.

# Experiments and Analysis

## Is it Possible for animat with restricted boltzmann machine based controller to perform *online* learning?

rbm is a form of "example-based" learning. WIth a complete dataset, rbm can learn the "pattern" efficiently. However, complete dataset is often not available in animat training, so animat needs to learn on-the-fly.

In this experiment, animat has access to a oracle that tells whether animat's output to a given input is good or bad. The rbm is pretrained to give "common sense" of what it should do. The number of pretrain example is 1000, while the number of possible inputs are exponential (2^25). Therefore, the pretrain examples only give a very rough energy curve for rbm.

The experiment has 1000 iterations, and in each iteration, 8 animats with different skin colors *daydream* and produce the outputs.
For each animat, its output is fed into the oracle. If the output is good, the output is added to a set for later training/learning. Due to memory size limitation and more realistic approximation to nature, the set has a fixed size. The expected behavior is defined by the oracle: food finding while attacking on hostile enemy on high food grid.

The result shows that sufficient examples are needed for pretraining. In this experiment, the minimum pretrain example for input space of 2^25 is around 400. After pretraining, animat can refine its energy curve with the online examples selected by the oracle. It is noteworthy that pretraining+online learning is far more efficient than pretraining alone. This is because pretraining inputs are randomly generated, and its distribution is distinguishable from the distribution the animat actually preceives. As a result, someof the examples have very little probability to be useful; online learning on the other hand generates very "valuable" examples.

Pure online learning is surprisingly very inefficient. This might result from that on the energy curve where there is no example, the output is very random. As the number of output increases, it is very hard to obrain a good/desired output oracle; theoratically, expected learning iteration of this online learning algorithm exponentially increases as output increases linearly. One possible solution is to define a quality score which extends oracle's binary decision to a continous space, and that work is not done in this phase.

*Pretain+online learning of food finding and fight over resource*
![base](http://giant.gfycat.com/HeavyMellowGuineafowl.gif)

*Added animat relation and skin color preference*
```
There are 3 skin colors: white, gray and black.
Each animat has opinion on skin color; if the other has different skin color and the impression on that skin color is bad, there is a higher chance of attacking.
```
![extend](http://giant.gfycat.com/NeedyQuerulousAmericancrayfish.gif)

## Whether mini-batch in rbm helps animat learning

In pratical guide to training restricted boltzmann machine, rbm training uses "mini-batch". We conduct an experiment where the mini-batch size is 20 vs no mini-batch. The result shows that mini-batch is crucial in rbm training. The one without mini-batch is very ineffective.  According to Hinton's guide, "Increasing the mini-batch size by a factor of N leads to a more reliable gradient estimate but it does not increase the maximum stable learning rate by a factor of N, so the net effect is that the weight updates are smaller per gradient evaluation." This experiment validates such expectation.

As a side note, setting good example size to a certain value in the previous experiment does not hurt learning. because we can see this as our "mini-batch". 

*Not using mini-batch*
![no-batch](http://giant.gfycat.com/TidyScarceEmu.gif)

## How to approximate input output with rbm

While rbm is able to learn by patterns, it is often unclear to beginners how it can be used to approximate functions with inputs and outputs. Here are the precise steps for training inputs and outputs that I work out(not documented anywhere):

1. concatenate input and output binary vector in order as visible layer. The concatenated input and output vector is clamped as training examples, and then run contrastive divergence (in this experiment cd-10) as ususal to get proper weights for both input and output.

2. in the reconstruction, only input is given, and output is left to be computed. Create the vector that is has length equal to input+output and fill input part. For each contrastive divergence step, reset input part of the visible vector. After some iteration, the output is the remaining part of the input vector. 

Note that weight calculation and reconstruction can use cd-k where k>=1. But k=10 is commonly used as a more precise approximation of gibbs sampling than k=1.

*Before using input clamping reconstruction*
![before](http://giant.gfycat.com/VainWastefulAmericanwigeon.gif)


*After using input clamping in reconstruction*
![after](http://giant.gfycat.com/ContentOptimisticAardwolf.gif)

# References

* A practical guide to training restricted boltzmann machines: https://www.cs.toronto.edu/~hinton/absps/guideTR.pdf
* Deep learning rbm tutorial: http://deeplearning.net/tutorial/rbm.html
* A small collection of neural network algorithms in Go: https://github.com/r9y9/nnet
* Restricted Boltzmann Machines in Python: https://github.com/echen/restricted-boltzmann-machines
