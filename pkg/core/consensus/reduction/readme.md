## Binary Reduction

The Binary Reduction algorithm lays at the core of SBA\*. It converts the problem of reaching consensus on arbitrary values to reaching consensus on one of two values. It is an adaptation of the Turpin and Coan algorithm, originally concocted to solve the general Byzantine agreement when given a binary Byzantine agreement algorithm as a subroutine, for `n > 3f` (with `n` defined as total number of nodes and `f` defined as adversarial nodes).

Unlike other implementations, which normally utilize the original algorithm, Binary Reduction adopted in SBA\* follows a two-step approach, with the input of the second step depending on the output of the first one.

If no consensus have been reached on a uniform value, the algorithm returns a default value and waits for the next instantiation.

Binary Reduction acts as a uniform value extraction function which is then fed through the Block Agreement algorithm before exiting the loop in case of a successful termination of the Block Agreement algorithm.

### Values

#### Block Reduction Event

| Field           | Type    |
| --------------- | ------- |
| pubkeyBLS       | BLS Signature |
| round           | uint64  |
| step            | uint64  |
| blockhash       | uint256 |
| signedblockhash | BLS Signature |

### API

- Launch(eventbus, committee, duration) - Launches the reduction component, which publishes on the `Reduction` and `ReductionResultTopic` topics.

### Architecture

The `Block Reducer` component follows the event driven paradigm. It is connected to the node's `EventBus` through a `broker` and it delegates event-specific operations to its `EventHandler` implementation.

The `Reducer` entity is generic and spawns two `eventStopWatch` (one per step) to regulate the collection of the events and handle eventual timeout.

Like all the other consensus components, collection of the events and their marshalling/unmarshalling is delegated to an `EventFilter`.

#### Block Reducer Diagram

![](docs/Block%20Reduction.jpg)
