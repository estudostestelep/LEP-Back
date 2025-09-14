import { useEffect, useState } from "react";
import { productService } from "../../api/productService";
import { orderService } from "../../api/ordersService";

export default function PublicCardapio() {
  const [products, setProducts] = useState<any[]>([]);
  const [mesa, setMesa] = useState("");
  const [pedido, setPedido] = useState<{ itemId: string, quantidade: number }[]>([]);

  useEffect(() => {
    loadProducts();
  }, []);

  const loadProducts = async () => {
    const data = await productService.getAll();
    setProducts(data.data);
  };

  const addToPedido = (itemId: string) => {
    setPedido((prev) => [...prev, { itemId, quantidade: 1 }]);
  };

  const handleSendPedido = async () => {
    if (!mesa) {
      alert("Informe o número da mesa!");
      return;
    }
    await orderService.create({
      mesaId: mesa,
      items: pedido.map(item => ({
        produtoId: item.itemId,
        quantity: item.quantidade
      }))
    });
    alert("Pedido enviado com sucesso!");
    setPedido([]);
  };

  return (
    <div className="p-4">
      <h1 className="text-xl mb-4">Cardápio</h1>
      <input
        placeholder="Número da mesa"
        value={mesa}
        onChange={(e) => setMesa(e.target.value)}
        className="border p-2 mb-4"
      />
      <div className="grid grid-cols-2 gap-4">
        {products.map((p) => (
          <div key={p.id} className="border p-3 rounded">
            <h2 className="font-bold">{p.nome}</h2>
            <p>{p.descricao}</p>
            <p className="font-semibold">R$ {p.preco}</p>
            <button
              onClick={() => addToPedido(p.id)}
              className="bg-green-500 text-white px-3 py-1 mt-2 rounded"
            >
              Adicionar
            </button>
          </div>
        ))}
      </div>
      {pedido.length > 0 && (
        <button
          onClick={handleSendPedido}
          className="bg-blue-600 text-white px-4 py-2 mt-6 rounded"
        >
          Finalizar Pedido
        </button>
      )}
    </div>
  );
}
