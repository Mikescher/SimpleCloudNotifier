// This class is useful togther with FutureBuilder
// Unfortunately Future.value(x) in FutureBuilder always results in one frame were snapshot.connectionState is waiting
// Whit way we can set the ImmediateFuture.value directly and circumvent that.

class ImmediateFuture<T> {
  final Future<T> future;
  final T? value;

  ImmediateFuture(this.future, this.value);

  ImmediateFuture.ofFuture(Future<T> v)
      : future = v,
        value = null;

  ImmediateFuture.ofValue(T v)
      : future = Future.value(v),
        value = v;
}
