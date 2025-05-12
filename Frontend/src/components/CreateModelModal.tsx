import { useState } from 'react';

interface CreateModelModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (modelData: any) => void;
}

const CreateModelModal: React.FC<CreateModelModalProps> = ({
  isOpen,
  onClose,
  onSubmit,
}) => {
  const [formData, setFormData] = useState({
    numberRequestsUsed: 1000,
    percentTrainData: 80,
    percentNormalRequests: 70,
    numTrees: 100,
    maxDepth: 10,
    minSamplesLeaf: 2,
    minSamplesSplit: 5,
    maxFeatures: 'sqrt',
    criterion: 'gini',
  });

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: isNaN(Number(value)) ? value : Number(value),
    }));
  };

  const handleSubmit = () => {
    onSubmit(formData);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 overflow-auto">
      <div className="bg-white p-8 rounded-lg shadow-lg w-[600px] max-h-screen overflow-y-auto">
        <h2 className="text-xl font-semibold mb-4">Create AI Model</h2>
        <div className="grid grid-cols-2 gap-x-4 gap-y-4">
          <div>
            <label className="block mb-1">Number of Requests Used</label>
            <input
              type="number"
              name="numberRequestsUsed"
              value={formData.numberRequestsUsed}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded"
            />
          </div>

          <div>
            <label className="block mb-1">Percent Train Data</label>
            <input
              type="number"
              name="percentTrainData"
              value={formData.percentTrainData}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded"
            />
          </div>

          <div>
            <label className="block mb-1">Percent Normal Requests</label>
            <input
              type="number"
              name="percentNormalRequests"
              value={formData.percentNormalRequests}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded"
            />
          </div>

          <div>
            <label className="block mb-1">Number of Trees</label>
            <input
              type="number"
              name="numTrees"
              value={formData.numTrees}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded"
            />
          </div>

          <div>
            <label className="block mb-1">Max Depth</label>
            <input
              type="number"
              name="maxDepth"
              value={formData.maxDepth}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded"
            />
          </div>

          <div>
            <label className="block mb-1">Min Samples Split</label>
            <input
              type="number"
              name="minSamplesSplit"
              value={formData.minSamplesSplit}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded"
            />
          </div>

          <div>
            <label className="block mb-1">Min Samples Leaf</label>
            <input
              type="number"
              name="minSamplesLeaf"
              value={formData.minSamplesLeaf}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded"
            />
          </div>

          <div>
            <label className="block mb-1">Max Features</label>
            <select
              name="maxFeatures"
              value={formData.maxFeatures}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded"
            >
              <option value="sqrt">sqrt</option>
              <option value="log2">log2</option>
              <option value="none">none</option>
            </select>
          </div>

          <div className="col-span-2">
            <label className="block mb-1">Criterion</label>
            <select
              name="criterion"
              value={formData.criterion}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded"
            >
              <option value="gini">gini</option>
              <option value="entropy">entropy</option>
            </select>
          </div>
        </div>

        <div className="mt-6 flex justify-end gap-4">
          <button className="px-4 py-2 bg-gray-200 rounded" onClick={onClose}>
            Cancel
          </button>
          <button
            className="px-4 py-2 bg-blue-500 text-white rounded"
            onClick={handleSubmit}
          >
            Create Model
          </button>
        </div>
      </div>
    </div>
  );
};

export default CreateModelModal;
