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

items = {
  {"string", "hello"},
  {"integer", 42},
  {"float", 3.14},
  {"atom", :atom_key},
  {"list", [1, 2, 3]},
  {"list-large-value", [1, 2, 3, 4, 256 + 1]},
  {"map", %{"key" => "value"}},
  {"binary", <<1, 2, 3, 4>>},
  {"tuple", {:tuple, 42, 3.14}},
  {"complex-tuple", {
    {:tuple, 42, 3.14},
    %{"key" => "value"},
    %{:atom_key => 123},
    3.14,
    42,
    "hello",
    <<1, 2, 3, 4>>,
    [1, 2, 3, 4, 5]
  }}
}

Enum.each(Tuple.to_list(items), fn {type, value} ->
  binary = :erlang.term_to_binary(value)
  File.write!(Path.join("round-trip", "input_#{type}.bin"), binary)
end)
