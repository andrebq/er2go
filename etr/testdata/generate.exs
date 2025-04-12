term = {
  {:tuple, 42, 3.14},
  %{"key" => "value", :atom_key => 123},
  3.14,
  42,
  "hello",
  <<1, 2, 3, 4>>,
  [1, 2, 3, 4, 5]
}

binary = :erlang.term_to_binary(term)
File.write!("input.bin", binary)
