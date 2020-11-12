namespace Unifi.Models.Shared
{
    /// <summary>
    /// The shared resource across the application
    /// </summary>
    public record Resource
    {
        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Mac { get; init; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool State { get; init; } = false;

        /// <inheritdoc />
        public override string ToString() => $"Mac: {this.Mac}; State: {this.State}";
    }
}
