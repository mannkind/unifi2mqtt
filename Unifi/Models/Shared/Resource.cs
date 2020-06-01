namespace Unifi.Models.Shared
{
    /// <summary>
    /// The shared resource across the application
    /// </summary>
    public class Resource
    {
        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Mac { get; set; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool State { get; set; } = false;

        /// <inheritdoc />
        public override string ToString() => $"Mac: {this.Mac}; State: {this.State}";
    }
}
